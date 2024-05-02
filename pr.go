package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func branchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func getPrefixes() []string {
	prefixes := strings.Split(os.Getenv("JIRA_PREFIXES"), ",")

	for i, prefix := range prefixes {
		prefixes[i] = strings.TrimSpace(prefix)
	}

	return prefixes
}

func getDefaultTitle(prefixes []string, branch string) (string, string) {
	if len(prefixes) == 0 {
		return "", capitalize(branch)
	}

	re := regexp.MustCompile(`(?i)^(` + strings.Join(prefixes, "|") + `)-\d+`)
	prefix := re.FindString(branch)
	title := strings.TrimSpace(strings.ReplaceAll(re.ReplaceAllString(branch, ""), "-", " "))

	if len(title) > 0 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}

	return strings.ToUpper(prefix), title
}

type answers struct {
	Title       string
	Description string
}

func buildForm(answers *answers) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Value(&answers.Title),

			huh.NewText().
				Title("Description").
				CharLimit(400).
				Value(&answers.Description),
		),
	).WithKeyMap(&huh.KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c")),
		Input: huh.InputKeyMap{
			AcceptSuggestion: key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "complete")),
			Prev:             key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:             key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
			Submit:           key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		},
		Text: huh.TextKeyMap{
			Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
			Next:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
			Submit:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
			NewLine: key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "new line")),
			Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
		},
	})
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func main() {
	flags := os.Args[1:]

	branch := branchName()
	prefixes := getPrefixes()
	prefix, title := getDefaultTitle(prefixes, branch)

	answers := &answers{Title: title}
	form := buildForm(answers)

	err := form.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	finalTitle := answers.Title

	// Automatically add the ticket to the title if present
	if prefix != "" {
		finalTitle = prefix + " " + finalTitle
	}

	fmt.Println("Creating pull request...")
	fmt.Println("gh pr create", strings.Join(flags, " "), "--title", quote(finalTitle), "--body", quote(answers.Description))
	fmt.Println()

	// Push first to ensure the command will succeed
	cmd := exec.Command("git", "push")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to push branch %s\n", err)
	}

	args := []string{"pr", "create", "--title", finalTitle, "--body", answers.Description}
	args = append(args, flags...)

	cmd = exec.Command("gh", args...)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to create pull request %s\n", err)
	}

	fmt.Println("Pull request created!")
}
