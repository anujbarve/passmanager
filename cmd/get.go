// cmd/get.go
package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	getID   string
	getCopy bool
	getShow bool
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a credential",
	Run:   runGet,
}

func init() {
	getCmd.Flags().StringVarP(&getID, "id", "i", "", "Credential ID (required)")
	getCmd.Flags().BoolVarP(&getCopy, "copy", "c", false, "Copy password to clipboard")
	getCmd.Flags().BoolVarP(&getShow, "show", "s", false, "Show password in output")
	getCmd.MarkFlagRequired("id")
}

func runGet(cmd *cobra.Command, args []string) {
	_, client, cryptoSvc := authenticate()
	defer cryptoSvc.SecureClear()

	cred, err := client.GetCredential(getID)
	if err != nil {
		fmt.Printf("❌ Credential not found: %v\n", err)
		os.Exit(1)
	}

	// Decrypt password
	password, err := cryptoSvc.Decrypt(cred.EncryptedPassword)
	if err != nil {
		fmt.Printf("❌ Failed to decrypt password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n📋 Credential Details")
	fmt.Println("=====================")
	fmt.Printf("Title:    %s\n", cred.Title)
	fmt.Printf("Username: %s\n", cred.Username)
	fmt.Printf("URL:      %s\n", cred.URL)
	fmt.Printf("Category: %s\n", cred.Category)

	if getShow {
		fmt.Printf("Password: %s\n", password)
	} else {
		fmt.Printf("Password: %s\n", "********")
	}

	if cred.Notes != "" {
		notes, _ := cryptoSvc.Decrypt(cred.Notes)
		fmt.Printf("Notes:    %s\n", notes)
	}

	if getCopy {
		if err := clipboard.WriteAll(password); err != nil {
			fmt.Printf("❌ Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Println("\n✅ Password copied to clipboard!")
		}
	}
}
