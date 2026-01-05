// internal/ui/ui.go
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	AppName    = "PassManager"
	AppVersion = "1.0.0"
)

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintBanner() {
	banner := `
   ____                __  __                                   
  |  _ \ __ _ ___ ___ |  \/  | __ _ _ __   __ _  __ _  ___ _ __ 
  | |_) / _' / __/ __|| |\/| |/ _' | '_ \ / _' |/ _' |/ _ \ '__|
  |  __/ (_| \__ \__ \| |  | | (_| | | | | (_| | (_| |  __/ |   
  |_|   \__,_|___/___/|_|  |_|\__,_|_| |_|\__,_|\__, |\___|_|   
                                                |___/           
`
	fmt.Println(Cyan + banner + Reset)
	fmt.Printf("  %s%sSecure Password Manager v%s%s\n", Bold, White, AppVersion, Reset)
	fmt.Printf("  %s\n\n", Subtle("Your passwords, encrypted locally, stored securely."))
}

func PrintDivider() {
	fmt.Println(Subtle(strings.Repeat("─", 60)))
}

func PrintSection(title string) {
	fmt.Printf("\n%s%s %s %s%s\n", Bold, Cyan, "▶", title, Reset)
	PrintDivider()
}

func PromptContinue() {
	fmt.Printf("\n%s", Subtle("Press Enter to continue..."))
	fmt.Scanln()
}

func PrintKeyValue(key, value string) {
	fmt.Printf("  %s%-15s%s %s\n", Dim, key+":", Reset, value)
}

func PrintCredentialCard(id, title, username, url, category string, showPassword bool, password string) {
	fmt.Println()
	fmt.Printf("  %s┌─────────────────────────────────────────────────┐%s\n", Cyan, Reset)
	fmt.Printf("  %s│%s %s%-47s%s %s│%s\n", Cyan, Reset, Bold+White, truncate(title, 47), Reset, Cyan, Reset)
	fmt.Printf("  %s├─────────────────────────────────────────────────┤%s\n", Cyan, Reset)
	fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"ID:"+Reset, "", truncate(id, 33), "", Cyan, Reset)
	fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"Username:"+Reset, "", truncate(username, 33), "", Cyan, Reset)
	fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"URL:"+Reset, "", truncate(url, 33), "", Cyan, Reset)
	fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"Category:"+Reset, "", truncate(category, 33), "", Cyan, Reset)
	
	if showPassword {
		fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"Password:"+Reset, Green, truncate(password, 33), Reset, Cyan, Reset)
	} else {
		fmt.Printf("  %s│%s  %-12s %s%-33s%s %s│%s\n", Cyan, Reset, Dim+"Password:"+Reset, Yellow, "••••••••••••", Reset, Cyan, Reset)
	}
	
	fmt.Printf("  %s└─────────────────────────────────────────────────┘%s\n", Cyan, Reset)
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}