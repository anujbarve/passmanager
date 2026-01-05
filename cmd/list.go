// cmd/list.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listSearch string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all credentials",
	Run:   runList,
}

func init() {
	listCmd.Flags().StringVarP(&listSearch, "search", "s", "", "Search by title, username, or URL")
}

func runList(cmd *cobra.Command, args []string) {
	_, client, cryptoSvc := authenticate()
	defer cryptoSvc.SecureClear()

	creds, err := client.ListCredentials(listSearch)
	if err != nil {
		fmt.Printf("❌ Failed to list credentials: %v\n", err)
		os.Exit(1)
	}

	if len(creds) == 0 {
		fmt.Println("📭 No credentials found")
		return
	}

	fmt.Println("\n🔐 Stored Credentials")
	fmt.Println("=====================")
	fmt.Printf("%-20s %-25s %-30s %-15s\n", "ID", "TITLE", "USERNAME", "CATEGORY")
	fmt.Println("-------------------- ------------------------- ------------------------------ ---------------")

	for _, cred := range creds {
		title := truncate(cred.Title, 23)
		username := truncate(cred.Username, 28)
		fmt.Printf("%-20s %-25s %-30s %-15s\n", cred.ID, title, username, cred.Category)
	}

	fmt.Printf("\nTotal: %d credential(s)\n", len(creds))
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-2] + ".."
	}
	return s
}
