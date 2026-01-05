// cmd/generate.go
package cmd

import (
	"fmt"
	"os"

	"passmanager/internal/crypto"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	genLength  int
	genSymbols bool
	genCopy    bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a secure password",
	Run:   runGenerate,
}

func init() {
	generateCmd.Flags().IntVarP(&genLength, "length", "l", 20, "Password length")
	generateCmd.Flags().BoolVarP(&genSymbols, "symbols", "s", true, "Include symbols")
	generateCmd.Flags().BoolVarP(&genCopy, "copy", "c", false, "Copy to clipboard")
}

func runGenerate(cmd *cobra.Command, args []string) {
	password, err := crypto.GeneratePassword(genLength, genSymbols)
	if err != nil {
		fmt.Printf("❌ Failed to generate password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔑 Generated Password: %s\n", password)

	if genCopy {
		if err := clipboard.WriteAll(password); err != nil {
			fmt.Printf("❌ Failed to copy: %v\n", err)
		} else {
			fmt.Println("✅ Copied to clipboard!")
		}
	}
}
