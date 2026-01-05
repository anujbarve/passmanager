// main.go
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"passmanager/internal/config"
	"passmanager/internal/crypto"
	"passmanager/internal/database"
	"passmanager/internal/models"
	"passmanager/internal/session"
	"passmanager/internal/ui"

	"github.com/atotto/clipboard"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)


func main() {
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n" + ui.Warning("Shutting down securely..."))
		sess := session.GetSession()
		sess.Logout()
		os.Exit(0)
	}()

	// Start application
	ui.ClearScreen()
	ui.PrintBanner()

	// Check if initialized
	if !config.Exists() {
		runSetup()
	}

	// Main loop
	runMainLoop()
}

func runSetup() {
	ui.PrintSection("First Time Setup")

	fmt.Println(ui.Info("Let's set up your secure password vault.\n"))

	// Get PocketBase URL
	pbURL, err := ui.InputPrompt("PocketBase URL", "http://127.0.0.1:8090", validateURL)
	if err != nil {
		fmt.Println(ui.Error("Setup cancelled"))
		os.Exit(1)
	}

	// Test connection
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Testing connection..."
	s.Start()

	client := database.NewPocketBaseClient(pbURL)
	if err := client.TestConnection(); err != nil {
		s.Stop()
		fmt.Println(ui.Error(fmt.Sprintf("Cannot connect to PocketBase: %v", err)))
		fmt.Println(ui.Info("Make sure PocketBase is running: ./pocketbase serve"))
		os.Exit(1)
	}
	s.Stop()
	fmt.Println(ui.Success("Connected to PocketBase"))

	// Get admin credentials
	adminEmail, _ := ui.InputPrompt("Admin Email", "", validateEmail)
	adminPass, _ := ui.PasswordPrompt("Admin Password")

	s.Suffix = " Authenticating..."
	s.Start()

	if err := client.Authenticate(adminEmail, adminPass); err != nil {
		s.Stop()
		fmt.Println(ui.Error(fmt.Sprintf("Authentication failed: %v", err)))
		os.Exit(1)
	}
	s.Stop()
	fmt.Println(ui.Success("Authenticated successfully"))

	// Create master password
	fmt.Println()
	fmt.Println(ui.Info("Now create your master password."))
	fmt.Println(ui.Subtle("  This password encrypts all your data locally."))
	fmt.Println(ui.Subtle("  It cannot be recovered if lost!"))
	fmt.Println()

	var masterPass string
	for {
		masterPass, _ = ui.PasswordPrompt("Master Password (min 12 chars)")
		if len(masterPass) < 12 {
			fmt.Println(ui.Error("Password must be at least 12 characters"))
			continue
		}

		confirmPass, _ := ui.PasswordPrompt("Confirm Master Password")
		if masterPass != confirmPass {
			fmt.Println(ui.Error("Passwords don't match"))
			continue
		}
		break
	}

	// Generate salt and save config
	s.Suffix = " Setting up vault..."
	s.Start()

	salt, _ := crypto.GenerateSalt()
	passwordHash := crypto.HashMasterPassword(masterPass, salt)

	vaultConfig := models.VaultConfig{
		Salt:         base64.StdEncoding.EncodeToString(salt),
		PasswordHash: passwordHash,
	}

	if err := client.SaveVaultConfig(vaultConfig); err != nil {
		s.Stop()
		fmt.Println(ui.Error(fmt.Sprintf("Failed to save vault config: %v", err)))
		fmt.Println(ui.Info("Make sure 'vault_config' collection exists in PocketBase"))
		os.Exit(1)
	}

	cfg := &config.Config{
		PocketBaseURL: pbURL,
		AdminEmail:    adminEmail,
		Settings:      models.DefaultSettings(),
		Initialized:   true,
	}

	if err := cfg.Save(); err != nil {
		s.Stop()
		fmt.Println(ui.Error(fmt.Sprintf("Failed to save config: %v", err)))
		os.Exit(1)
	}

	s.Stop()
	fmt.Println(ui.Success("Vault created successfully!"))
	fmt.Println()
	fmt.Println(ui.Warning("IMPORTANT: Remember your master password!"))
	fmt.Println(ui.Subtle("  It cannot be recovered if lost."))
	ui.PromptContinue()
}

func runMainLoop() {
	sess := session.GetSession()

	for {
		ui.ClearScreen()
		ui.PrintBanner()

		// Check session
		if !sess.IsAuthenticated() {
			if !authenticate() {
				continue
			}
		}

		// Show session status
		remaining := sess.GetTimeRemaining()
		fmt.Printf("%s Session active (expires in %s)\n",
			ui.Subtle("🔓"),
			ui.Subtle(formatDuration(remaining)))

		// Show main menu
		choice, err := ui.MainMenu()
		if err != nil {
			if err == promptui.ErrInterrupt {
				handleExit()
			}
			continue
		}

		sess.UpdateActivity()

		switch choice {
		case "Add Credential":
			handleAddCredential()
		case "List Credentials":
			handleListCredentials()
		case "Search Credentials":
			handleSearchCredentials()
		case "Get Credential":
			handleGetCredential()
		case "Generate Password":
			handleGeneratePassword()
		case "Delete Credential":
			handleDeleteCredential()
		case "Change Master Password":
			handleChangeMasterPassword()
		case "Lock Vault":
			sess.Logout()
			fmt.Println(ui.Success("Vault locked"))
			ui.PromptContinue()
		case "Settings":
			handleSettings()
		case "Help":
			handleHelp()
		case "Exit":
			handleExit()
		}
	}
}

func authenticate() bool {
	ui.PrintSection("Unlock Vault")

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(ui.Error("Configuration not found. Please run setup."))
		return false
	}

	// Get admin password
	adminPass, err := ui.PasswordPrompt("Admin Password")
	if err != nil {
		return false
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Connecting..."
	s.Start()

	client := database.NewPocketBaseClient(cfg.PocketBaseURL)
	if err := client.Authenticate(cfg.AdminEmail, adminPass); err != nil {
		s.Stop()
		fmt.Println(ui.Error("Authentication failed"))
		ui.PromptContinue()
		return false
	}

	vaultConfig, err := client.GetVaultConfig()
	if err != nil {
		s.Stop()
		fmt.Println(ui.Error("Failed to load vault configuration"))
		ui.PromptContinue()
		return false
	}

	s.Stop()

	// Get master password
	masterPass, err := ui.PasswordPrompt("Master Password")
	if err != nil {
		return false
	}

	salt, _ := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	passwordHash := crypto.HashMasterPassword(masterPass, salt)

	if passwordHash != vaultConfig.PasswordHash {
		fmt.Println(ui.Error("Invalid master password"))
		ui.PromptContinue()
		return false
	}

	cryptoSvc := crypto.NewCryptoService(masterPass, salt)

	sess := session.GetSession()
	sess.Login(client, cryptoSvc, salt)
	sess.SetTimeout(time.Duration(cfg.Settings.SessionTimeout) * time.Minute)

	fmt.Println(ui.Success("Vault unlocked!"))
	time.Sleep(500 * time.Millisecond)

	return true
}

func handleAddCredential() {
	ui.ClearScreen()
	ui.PrintSection("Add New Credential")

	sess := session.GetSession()
	cfg, _ := config.Load()

	// Get credential details
	title, _ := ui.InputPrompt("Title", "", validateRequired)
	username, _ := ui.InputPrompt("Username/Email", "", nil)
	urlInput, _ := ui.InputPrompt("URL", "", nil)
	category, _ := ui.InputPrompt("Category", cfg.Settings.DefaultCategory, nil)
	notes, _ := ui.InputPrompt("Notes (optional)", "", nil)

	// Password options
	passOptions := []string{
		"🎲 Generate secure password",
		"✏️  Enter password manually",
	}
	_, passChoice, _ := ui.SelectFromList("Password", passOptions)

	var password string
	if strings.Contains(passChoice, "Generate") {
		length := cfg.Settings.PasswordLength
		lengthStr, _ := ui.InputPrompt("Password length", strconv.Itoa(length), validateNumber)
		length, _ = strconv.Atoi(lengthStr)

		password, _ = crypto.GeneratePassword(length, cfg.Settings.IncludeSymbols)
		fmt.Printf("\n%s Generated: %s%s%s\n", ui.Subtle("🔑"), ui.Green+ui.Bold, password, ui.Reset)
	} else {
		password, _ = ui.PasswordPrompt("Password")
	}

	// Encrypt and save
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Encrypting and saving..."
	s.Start()

	cryptoSvc := sess.GetCrypto()
	encryptedPassword, _ := cryptoSvc.Encrypt(password)

	encryptedNotes := ""
	if notes != "" {
		encryptedNotes, _ = cryptoSvc.Encrypt(notes)
	}

	cred := models.Credential{
		Title:             title,
		Username:          username,
		EncryptedPassword: encryptedPassword,
		URL:               urlInput,
		Notes:             encryptedNotes,
		Category:          category,
	}

	created, err := sess.GetDB().CreateCredential(cred)
	s.Stop()

	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to save: %v", err)))
	} else {
		fmt.Println(ui.Success(fmt.Sprintf("Credential saved! ID: %s", created.ID)))

		if ui.ConfirmPrompt("Copy password to clipboard?") {
			clipboard.WriteAll(password)
			fmt.Println(ui.Success("Password copied!"))
		}
	}

	ui.PromptContinue()
}

func handleListCredentials() {
	ui.ClearScreen()
	ui.PrintSection("All Credentials")

	sess := session.GetSession()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Loading..."
	s.Start()

	creds, err := sess.GetDB().ListCredentials("")
	s.Stop()

	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to load credentials: %v", err)))
		ui.PromptContinue()
		return
	}

	if len(creds) == 0 {
		fmt.Println(ui.Info("No credentials stored yet."))
		fmt.Println(ui.Subtle("Use 'Add Credential' to store your first password."))
		ui.PromptContinue()
		return
	}

	printCredentialsTable(creds)

	fmt.Printf("\n%s\n", ui.Subtle(fmt.Sprintf("Total: %d credential(s)", len(creds))))
	ui.PromptContinue()
}

func handleSearchCredentials() {
	ui.ClearScreen()
	ui.PrintSection("Search Credentials")

	query, _ := ui.InputPrompt("Search term", "", validateRequired)

	sess := session.GetSession()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Searching..."
	s.Start()

	creds, err := sess.GetDB().ListCredentials(query)
	s.Stop()

	if err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Search failed: %v", err)))
		ui.PromptContinue()
		return
	}

	if len(creds) == 0 {
		fmt.Println(ui.Info(fmt.Sprintf("No credentials found matching '%s'", query)))
		ui.PromptContinue()
		return
	}

	printCredentialsTable(creds)

	fmt.Printf("\n%s\n", ui.Subtle(fmt.Sprintf("Found: %d credential(s)", len(creds))))

	// Option to view one
	if ui.ConfirmPrompt("View credential details?") {
		id, _ := ui.InputPrompt("Enter ID", "", validateRequired)
		viewCredential(id)
	}

	ui.PromptContinue()
}

func handleGetCredential() {
	ui.ClearScreen()
	ui.PrintSection("Get Credential")

	id, _ := ui.InputPrompt("Credential ID", "", validateRequired)
	viewCredential(id)
	ui.PromptContinue()
}

func viewCredential(id string) {
	sess := session.GetSession()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Loading..."
	s.Start()

	cred, err := sess.GetDB().GetCredential(id)
	s.Stop()

	if err != nil {
		fmt.Println(ui.Error("Credential not found"))
		return
	}

	cryptoSvc := sess.GetCrypto()
	password, _ := cryptoSvc.Decrypt(cred.EncryptedPassword)

	notes := ""
	if cred.Notes != "" {
		notes, _ = cryptoSvc.Decrypt(cred.Notes)
	}

	ui.PrintCredentialCard(cred.ID, cred.Title, cred.Username, cred.URL, cred.Category, false, "")

	if notes != "" {
		fmt.Printf("\n  %sNotes:%s %s\n", ui.Dim, ui.Reset, notes)
	}

	// Actions menu
	actions := []string{
		"👁️  Show password",
		"📋 Copy password to clipboard",
		"📋 Copy username to clipboard",
		"🔙 Go back",
	}

	for {
		_, action, _ := ui.SelectFromList("Action", actions)

		switch {
		case strings.Contains(action, "Show password"):
			fmt.Printf("\n  %sPassword:%s %s%s%s\n", ui.Dim, ui.Reset, ui.Green, password, ui.Reset)
		case strings.Contains(action, "Copy password"):
			clipboard.WriteAll(password)
			fmt.Println(ui.Success("Password copied to clipboard!"))
		case strings.Contains(action, "Copy username"):
			clipboard.WriteAll(cred.Username)
			fmt.Println(ui.Success("Username copied to clipboard!"))
		case strings.Contains(action, "Go back"):
			return
		}
	}
}

func handleGeneratePassword() {
	ui.ClearScreen()
	ui.PrintSection("Generate Password")

	cfg, _ := config.Load()

	lengthStr, _ := ui.InputPrompt("Password length", strconv.Itoa(cfg.Settings.PasswordLength), validateNumber)
	length, _ := strconv.Atoi(lengthStr)

	includeSymbols := ui.ConfirmPrompt("Include symbols (!@#$%...)?")

	password, err := crypto.GeneratePassword(length, includeSymbols)
	if err != nil {
		fmt.Println(ui.Error("Failed to generate password"))
		ui.PromptContinue()
		return
	}

	fmt.Println()
	fmt.Printf("  %s┌─────────────────────────────────────────────────┐%s\n", ui.Cyan, ui.Reset)
	fmt.Printf("  %s│%s  Generated Password:                            %s│%s\n", ui.Cyan, ui.Reset, ui.Cyan, ui.Reset)
	fmt.Printf("  %s│%s  %s%-45s%s %s│%s\n", ui.Cyan, ui.Reset, ui.Green+ui.Bold, password, ui.Reset, ui.Cyan, ui.Reset)
	fmt.Printf("  %s└─────────────────────────────────────────────────┘%s\n", ui.Cyan, ui.Reset)
	fmt.Println()

	if ui.ConfirmPrompt("Copy to clipboard?") {
		clipboard.WriteAll(password)
		fmt.Println(ui.Success("Password copied!"))
	}

	ui.PromptContinue()
}

func handleDeleteCredential() {
	ui.ClearScreen()
	ui.PrintSection("Delete Credential")

	id, _ := ui.InputPrompt("Credential ID to delete", "", validateRequired)

	sess := session.GetSession()

	// Show credential first
	cred, err := sess.GetDB().GetCredential(id)
	if err != nil {
		fmt.Println(ui.Error("Credential not found"))
		ui.PromptContinue()
		return
	}

	fmt.Printf("\n%s You are about to delete:\n", ui.Warning(""))
	fmt.Printf("  Title: %s%s%s\n", ui.Bold, cred.Title, ui.Reset)
	fmt.Printf("  Username: %s\n", cred.Username)
	fmt.Println()

	if !ui.ConfirmPrompt("Are you sure? This cannot be undone") {
		fmt.Println(ui.Info("Deletion cancelled"))
		ui.PromptContinue()
		return
	}

	// Double confirm
	confirm, _ := ui.InputPrompt("Type 'DELETE' to confirm", "", nil)
	if confirm != "DELETE" {
		fmt.Println(ui.Info("Deletion cancelled"))
		ui.PromptContinue()
		return
	}

	if err := sess.GetDB().DeleteCredential(id); err != nil {
		fmt.Println(ui.Error(fmt.Sprintf("Failed to delete: %v", err)))
	} else {
		fmt.Println(ui.Success("Credential deleted"))
	}

	ui.PromptContinue()
}

func handleChangeMasterPassword() {
	ui.ClearScreen()
	ui.PrintSection("Change Master Password")

	fmt.Println(ui.Warning("This will re-encrypt all your credentials."))
	fmt.Println(ui.Subtle("Make sure you have a backup before proceeding."))
	fmt.Println()

	if !ui.ConfirmPrompt("Continue?") {
		return
	}

	sess := session.GetSession()
	cfg, _ := config.Load()

	// Verify current password
	currentPass, _ := ui.PasswordPrompt("Current Master Password")
	currentHash := crypto.HashMasterPassword(currentPass, sess.GetSalt())

	vaultConfig, _ := sess.GetDB().GetVaultConfig()
	if currentHash != vaultConfig.PasswordHash {
		fmt.Println(ui.Error("Invalid current password"))
		ui.PromptContinue()
		return
	}

	// Get new password
	var newPass string
	for {
		newPass, _ = ui.PasswordPrompt("New Master Password (min 12 chars)")
		if len(newPass) < 12 {
			fmt.Println(ui.Error("Password must be at least 12 characters"))
			continue
		}

		confirmPass, _ := ui.PasswordPrompt("Confirm New Password")
		if newPass != confirmPass {
			fmt.Println(ui.Error("Passwords don't match"))
			continue
		}
		break
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Re-encrypting all credentials..."
	s.Start()

	// Get all credentials
	creds, _ := sess.GetDB().ListCredentials("")

	// Create new crypto service
	newSalt, _ := crypto.GenerateSalt()
	newCryptoSvc := crypto.NewCryptoService(newPass, newSalt)
	oldCryptoSvc := sess.GetCrypto()

	// Re-encrypt all credentials
	for _, cred := range creds {
		// Decrypt with old key
		password, _ := oldCryptoSvc.Decrypt(cred.EncryptedPassword)
		notes := ""
		if cred.Notes != "" {
			notes, _ = oldCryptoSvc.Decrypt(cred.Notes)
		}

		// Encrypt with new key
		newEncPassword, _ := newCryptoSvc.Encrypt(password)
		newEncNotes := ""
		if notes != "" {
			newEncNotes, _ = newCryptoSvc.Encrypt(notes)
		}

		// Update credential
		cred.EncryptedPassword = newEncPassword
		cred.Notes = newEncNotes
		sess.GetDB().UpdateCredential(cred.ID, cred)
	}

	// Update vault config
	newPasswordHash := crypto.HashMasterPassword(newPass, newSalt)
	vaultConfig.Salt = base64.StdEncoding.EncodeToString(newSalt)
	vaultConfig.PasswordHash = newPasswordHash
	sess.GetDB().UpdateVaultConfig(vaultConfig.ID, *vaultConfig)

	s.Stop()

	// Update session
	sess.Login(sess.GetDB(), newCryptoSvc, newSalt)
	sess.SetTimeout(time.Duration(cfg.Settings.SessionTimeout) * time.Minute)

	fmt.Println(ui.Success("Master password changed successfully!"))
	fmt.Println(ui.Warning("Remember your new password!"))

	ui.PromptContinue()
}

func handleSettings() {
	ui.ClearScreen()
	ui.PrintSection("Settings")

	cfg, _ := config.Load()

	for {
		fmt.Println()
		fmt.Printf("  %s1.%s Session Timeout: %s%d minutes%s\n",
			ui.Cyan, ui.Reset, ui.Bold, cfg.Settings.SessionTimeout, ui.Reset)
		fmt.Printf("  %s2.%s Clipboard Timeout: %s%d seconds%s\n",
			ui.Cyan, ui.Reset, ui.Bold, cfg.Settings.ClipboardTimeout, ui.Reset)
		fmt.Printf("  %s3.%s Default Category: %s%s%s\n",
			ui.Cyan, ui.Reset, ui.Bold, cfg.Settings.DefaultCategory, ui.Reset)
		fmt.Printf("  %s4.%s Default Password Length: %s%d%s\n",
			ui.Cyan, ui.Reset, ui.Bold, cfg.Settings.PasswordLength, ui.Reset)
		fmt.Printf("  %s5.%s Include Symbols by Default: %s%v%s\n",
			ui.Cyan, ui.Reset, ui.Bold, cfg.Settings.IncludeSymbols, ui.Reset)
		fmt.Printf("  %s6.%s Back to Main Menu\n", ui.Cyan, ui.Reset)
		fmt.Println()

		choice, _ := ui.InputPrompt("Select option (1-6)", "", nil)

		switch choice {
		case "1":
			val, _ := ui.InputPrompt("Session timeout (minutes)", strconv.Itoa(cfg.Settings.SessionTimeout), validateNumber)
			cfg.Settings.SessionTimeout, _ = strconv.Atoi(val)
			sess := session.GetSession()
			sess.SetTimeout(time.Duration(cfg.Settings.SessionTimeout) * time.Minute)
		case "2":
			val, _ := ui.InputPrompt("Clipboard timeout (seconds)", strconv.Itoa(cfg.Settings.ClipboardTimeout), validateNumber)
			cfg.Settings.ClipboardTimeout, _ = strconv.Atoi(val)
		case "3":
			cfg.Settings.DefaultCategory, _ = ui.InputPrompt("Default category", cfg.Settings.DefaultCategory, nil)
		case "4":
			val, _ := ui.InputPrompt("Default password length", strconv.Itoa(cfg.Settings.PasswordLength), validateNumber)
			cfg.Settings.PasswordLength, _ = strconv.Atoi(val)
		case "5":
			cfg.Settings.IncludeSymbols = ui.ConfirmPrompt("Include symbols by default?")
		case "6":
			cfg.Save()
			return
		}

		cfg.Save()
		fmt.Println(ui.Success("Setting updated"))
	}
}

func handleHelp() {
	ui.ClearScreen()
	ui.PrintSection("Help")

	helpText := `
  %s🔐 About PassManager%s
  PassManager is a secure, local-first password manager.
  All passwords are encrypted using AES-256-GCM before being stored.

  %s📋 Features:%s
  • Store unlimited passwords securely
  • Generate cryptographically secure passwords
  • Search and organize by categories
  • Copy passwords to clipboard
  • Session timeout for security

  %s🔒 Security:%s
  • Master password never stored
  • Argon2id key derivation (memory-hard)
  • AES-256-GCM authenticated encryption
  • Data encrypted locally before sending to server

  %s⌨️  Keyboard Shortcuts:%s
  • Ctrl+C: Lock vault and exit
  • Enter: Confirm selection
  • Type to filter in menus

  %s📖 Tips:%s
  • Use a strong master password (16+ characters)
  • Enable clipboard timeout in settings
  • Regularly backup your PocketBase data
  • Lock vault when stepping away

  %s🆘 Support:%s
  • GitHub: github.com/yourusername/passmanager
  • Email: support@example.com
`
	fmt.Printf(helpText,
		ui.Bold+ui.Cyan, ui.Reset,
		ui.Bold+ui.Cyan, ui.Reset,
		ui.Bold+ui.Cyan, ui.Reset,
		ui.Bold+ui.Cyan, ui.Reset,
		ui.Bold+ui.Cyan, ui.Reset,
		ui.Bold+ui.Cyan, ui.Reset,
	)

	ui.PromptContinue()
}

func handleExit() {
	fmt.Println()
	if ui.ConfirmPrompt("Exit PassManager?") {
		sess := session.GetSession()
		sess.Logout()
		fmt.Println(ui.Success("Vault locked. Goodbye! 👋"))
		os.Exit(0)
	}
}

// Helper functions

func printCredentialsTable(creds []models.Credential) {
	if len(creds) == 0 {
		fmt.Println(ui.Info("No credentials found"))
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Set header
	t.AppendHeader(table.Row{"ID", "Title", "Username", "Category"})

	// Add rows
	for _, cred := range creds {
		t.AppendRow(table.Row{
			truncateStr(cred.ID, 15),
			truncateStr(cred.Title, 25),
			truncateStr(cred.Username, 30),
			truncateStr(cred.Category, 15),
		})
	}

	// Style configuration
	t.SetStyle(table.Style{
		Name: "PassManager",
		Box: table.BoxStyle{
			BottomLeft:       "└",
			BottomRight:      "┘",
			BottomSeparator:  "┴",
			Left:             "│",
			LeftSeparator:    "├",
			MiddleHorizontal: "─",
			MiddleSeparator:  "┼",
			MiddleVertical:   "│",
			PaddingLeft:      " ",
			PaddingRight:     " ",
			Right:            "│",
			RightSeparator:   "┤",
			TopLeft:          "┌",
			TopRight:         "┐",
			TopSeparator:     "┬",
			UnfinishedRow:    "...",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.FgCyan, text.Bold},
			Row:    text.Colors{text.FgWhite},
			Footer: text.Colors{text.FgCyan},
		},
		Format: table.FormatOptions{
			Header: text.FormatDefault,
			Row:    text.FormatDefault,
			Footer: text.FormatDefault,
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateFooter:  false,
			SeparateHeader:  true,
			SeparateRows:    false,
		},
	})

	t.Render()
}

func truncateStr(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	return fmt.Sprintf("%dm %ds", m, s)
}

// Validators

func validateRequired(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("this field is required")
	}
	return nil
}

// Validators (continued in main.go)

func validateURL(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf("URL is required")
	}
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}

func validateEmail(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(input, "@") || !strings.Contains(input, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validateNumber(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf("number is required")
	}
	num, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("must be a valid number")
	}
	if num < 1 {
		return fmt.Errorf("must be greater than 0")
	}
	return nil
}