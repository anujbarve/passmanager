// internal/ui/menu.go
package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type MenuItem struct {
	Name        string
	Description string
	Icon        string
}

func (m MenuItem) String() string {
	return fmt.Sprintf("%s  %s", m.Icon, m.Name)
}

func MainMenu() (string, error) {
	items := []MenuItem{
		{Name: "Add Credential", Description: "Store a new password", Icon: "➕"},
		{Name: "List Credentials", Description: "View all stored passwords", Icon: "📋"},
		{Name: "Search Credentials", Description: "Find a specific password", Icon: "🔍"},
		{Name: "Get Credential", Description: "Retrieve a password by ID", Icon: "🔑"},
		{Name: "Generate Password", Description: "Create a secure password", Icon: "🎲"},
		{Name: "Delete Credential", Description: "Remove a stored password", Icon: "🗑️ "},
		{Name: "Change Master Password", Description: "Update your master password", Icon: "🔐"},
		{Name: "Lock Vault", Description: "Lock and require re-authentication", Icon: "🔒"},
		{Name: "Settings", Description: "Configure application settings", Icon: "⚙️ "},
		{Name: "Help", Description: "Show help information", Icon: "❓"},
		{Name: "Exit", Description: "Close the application", Icon: "🚪"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   fmt.Sprintf("%s {{ .Icon }}  {{ .Name | cyan | bold }} %s{{ .Description | faint }}%s", "▸", "(", ")"),
		Inactive: "  {{ .Icon }}  {{ .Name | white }} {{ .Description | faint }}",
		Selected: fmt.Sprintf("%s {{ .Icon }}  {{ .Name | green | bold }}", "✔"),
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("\n%s%s Main Menu %s", Bold+Cyan, "🔐", Reset),
		Items:     items,
		Templates: templates,
		Size:      11,
		HideHelp:  true,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[idx].Name, nil
}

func ConfirmPrompt(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return result == "y" || result == "Y"
}

func SelectFromList(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  10,
	}

	return prompt.Run()
}

func InputPrompt(label string, defaultVal string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultVal,
		Validate: validate,
	}

	return prompt.Run()
}

func PasswordPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '•',
	}

	return prompt.Run()
}