// cmd/delete.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var deleteID string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a credential",
	Run:   runDelete,
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteID, "id", "i", "", "Credential ID to delete (required)")
	deleteCmd.MarkFlagRequired("id")
}

func runDelete(cmd *cobra.Command, args []string) {
	_, client, cryptoSvc := authenticate()
	defer cryptoSvc.SecureClear()

	// Confirm deletion
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("⚠️  Are you sure you want to delete credential %s? (yes/no): ", deleteID)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "yes" {
		fmt.Println("❌ Deletion cancelled")
		return
	}

	if err := client.DeleteCredential(deleteID); err != nil {
		fmt.Printf("❌ Failed to delete credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Credential deleted successfully")
}
