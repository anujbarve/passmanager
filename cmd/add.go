// cmd/add.go
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"syscall"

	"passmanager/internal/config"
	"passmanager/internal/crypto"
	"passmanager/internal/database"
	"passmanager/internal/models"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	addTitle    string
	addUsername string
	addPassword string
	addURL      string
	addNotes    string
	addCategory string
	addGenerate bool
	addLength   int
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new credential",
	Run:   runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "Title/name for the credential (required)")
	addCmd.Flags().StringVarP(&addUsername, "username", "u", "", "Username/email")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "Password (will prompt if not provided)")
	addCmd.Flags().StringVarP(&addURL, "url", "l", "", "Website URL")
	addCmd.Flags().StringVarP(&addNotes, "notes", "n", "", "Additional notes")
	addCmd.Flags().StringVarP(&addCategory, "category", "c", "general", "Category")
	addCmd.Flags().BoolVarP(&addGenerate, "generate", "g", false, "Generate a random password")
	addCmd.Flags().IntVar(&addLength, "length", 20, "Generated password length")
	addCmd.MarkFlagRequired("title")
}

func runAdd(cmd *cobra.Command, args []string) {
	cfg, client, cryptoSvc := authenticate()
	defer cryptoSvc.SecureClear()

	// Handle password
	var password string
	if addGenerate {
		var err error
		password, err = crypto.GeneratePassword(addLength, true)
		if err != nil {
			fmt.Printf("❌ Failed to generate password: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🔑 Generated password: %s\n", password)
	} else if addPassword == "" {
		fmt.Print("Password: ")
		passBytes, _ := term.ReadPassword(int(syscall.Stdin))
		password = string(passBytes)
		fmt.Println()
	} else {
		password = addPassword
	}

	// Encrypt password
	encryptedPassword, err := cryptoSvc.Encrypt(password)
	if err != nil {
		fmt.Printf("❌ Failed to encrypt password: %v\n", err)
		os.Exit(1)
	}

	// Encrypt notes if provided
	encryptedNotes := ""
	if addNotes != "" {
		encryptedNotes, _ = cryptoSvc.Encrypt(addNotes)
	}

	// Create credential
	cred := models.Credential{
		Title:             addTitle,
		Username:          addUsername,
		EncryptedPassword: encryptedPassword,
		URL:               addURL,
		Notes:             encryptedNotes,
		Category:          addCategory,
	}

	_ = cfg // Use config if needed
	created, err := client.CreateCredential(cred)
	if err != nil {
		fmt.Printf("❌ Failed to save credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Credential saved successfully (ID: %s)\n", created.ID)
}

func authenticate() (*config.Config, *database.PocketBaseClient, *crypto.CryptoService) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("❌ Vault not initialized. Run 'passmanager init' first.")
		os.Exit(1)
	}

	// Get admin password
	fmt.Print("Admin Password: ")
	adminPassBytes, _ := term.ReadPassword(int(syscall.Stdin))
	adminPass := string(adminPassBytes)
	fmt.Println()

	// Connect to PocketBase
	client := database.NewPocketBaseClient(cfg.PocketBaseURL)
	if err := client.Authenticate(cfg.AdminEmail, adminPass); err != nil {
		fmt.Printf("❌ Authentication failed: %v\n", err)
		os.Exit(1)
	}

	// Get vault config
	vaultConfig, err := client.GetVaultConfig()
	if err != nil {
		fmt.Printf("❌ Failed to get vault config: %v\n", err)
		os.Exit(1)
	}

	// Get master password
	fmt.Print("Master Password: ")
	masterPassBytes, _ := term.ReadPassword(int(syscall.Stdin))
	masterPass := string(masterPassBytes)
	fmt.Println()

	// Verify master password
	salt, _ := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	passwordHash := crypto.HashMasterPassword(masterPass, salt)

	if passwordHash != vaultConfig.PasswordHash {
		fmt.Println("❌ Invalid master password")
		os.Exit(1)
	}

	cryptoSvc := crypto.NewCryptoService(masterPass, salt)

	return cfg, client, cryptoSvc
}
