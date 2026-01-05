// cmd/init.go
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"syscall"

	"passmanager/internal/config"
	"passmanager/internal/crypto"
	"passmanager/internal/database"
	"passmanager/internal/models"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the password vault",
	Long:  "Set up the password manager with master password and PocketBase connection",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🔐 Password Manager Setup")
	fmt.Println("========================")

	// Get PocketBase URL
	var pbURL string
	fmt.Print("PocketBase URL (e.g., http://127.0.0.1:8090): ")
	fmt.Scanln(&pbURL)
	
	// Clean URL
	pbURL = strings.TrimSpace(pbURL)
	pbURL = strings.TrimSuffix(pbURL, "/")

	// Create client and test connection
	client := database.NewPocketBaseClient(pbURL)
	
	fmt.Println("\n🔍 Testing connection to PocketBase...")
	if err := client.TestConnection(); err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Println("\n💡 Make sure PocketBase is running:")
		fmt.Println("   ./pocketbase serve")
		os.Exit(1)
	}
	fmt.Println("✅ PocketBase is reachable")

	// Get admin credentials
	var adminEmail string
	fmt.Print("\nAdmin/Superuser Email: ")
	fmt.Scanln(&adminEmail)
	adminEmail = strings.TrimSpace(adminEmail)

	fmt.Print("Admin/Superuser Password: ")
	adminPassBytes, _ := term.ReadPassword(int(syscall.Stdin))
	adminPass := string(adminPassBytes)
	fmt.Println()

	// Authenticate
	fmt.Println("\n🔐 Authenticating...")
	if err := client.Authenticate(adminEmail, adminPass); err != nil {
		fmt.Printf("❌ Failed to authenticate: %v\n", err)
		fmt.Println("\n💡 Troubleshooting:")
		fmt.Println("   1. Make sure you've created an admin/superuser in PocketBase")
		fmt.Println("   2. Go to PocketBase Admin UI → Settings → Admins")
		fmt.Println("   3. Or create a 'users' collection with email/password auth")
		os.Exit(1)
	}

	// Set master password
	fmt.Print("\nCreate Master Password (min 12 chars): ")
	masterPassBytes, _ := term.ReadPassword(int(syscall.Stdin))
	masterPass := string(masterPassBytes)
	fmt.Println()

	fmt.Print("Confirm Master Password: ")
	confirmPassBytes, _ := term.ReadPassword(int(syscall.Stdin))
	confirmPass := string(confirmPassBytes)
	fmt.Println()

	if masterPass != confirmPass {
		fmt.Println("❌ Passwords don't match")
		os.Exit(1)
	}

	if len(masterPass) < 12 {
		fmt.Println("❌ Master password must be at least 12 characters")
		os.Exit(1)
	}

	// Generate salt and hash
	salt, err := crypto.GenerateSalt()
	if err != nil {
		fmt.Printf("❌ Failed to generate salt: %v\n", err)
		os.Exit(1)
	}

	passwordHash := crypto.HashMasterPassword(masterPass, salt)

	// Save vault config to PocketBase
	vaultConfig := models.VaultConfig{
		Salt:         base64.StdEncoding.EncodeToString(salt),
		PasswordHash: passwordHash,
	}

	fmt.Println("\n📦 Saving vault configuration...")
	if err := client.SaveVaultConfig(vaultConfig); err != nil {
		fmt.Printf("❌ Failed to save vault config: %v\n", err)
		fmt.Println("\n💡 Make sure you've created the 'vault_config' collection:")
		fmt.Println("   - Go to PocketBase Admin UI")
		fmt.Println("   - Create collection 'vault_config'")
		fmt.Println("   - Add fields: salt (text), password_hash (text)")
		os.Exit(1)
	}

	// Save local config
	cfg := &config.Config{
		PocketBaseURL:  pbURL,
		AdminEmail:     adminEmail,
		SessionTimeout: 5,
	}

	if err := cfg.Save(); err != nil {
		fmt.Printf("❌ Failed to save local config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Password vault initialized successfully!")
	fmt.Println("⚠️  Remember your master password - it cannot be recovered!")
}
