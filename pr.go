package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func defaultBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func getDefaultTitle() string {
	cmd := exec.Command("git", "merge-base", "HEAD", defaultBranch())
	output, err := cmd.Output()

	if err != nil {
		return ""
	}

	cmd = exec.Command(
		"git", "log",
		"--reverse",
		"--format=%s", strings.TrimSpace(string(output))+"..HEAD",
	)
	output, err = cmd.Output()

	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		return lines[0]
	}

	return ""
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

func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	_, err = in.Write([]byte(text))
	if err != nil {
		return err
	}

	if err := in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}

func main() {
	flags := os.Args[1:]
	title := getDefaultTitle()

	answers := &answers{Title: title}
	form := buildForm(answers)
	if err := form.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Creating pull request...")
	printArgs := []string{"--title", quote(answers.Title), "--body", quote(answers.Description)}
	if len(flags) > 0 {
		fmt.Println("gh pr create", strings.Join(flags, " "), strings.Join(printArgs, " "))
	} else {
		fmt.Println("gh pr create", strings.Join(printArgs, " "))
	}

	// Push first to ensure the command will succeed
	cmd := exec.Command("git", "push")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to push branch %s\n", err)
	}

	args := []string{"pr", "create", "--title", answers.Title, "--body", answers.Description}
	args = append(args, flags...)

	var stdout bytes.Buffer
	cmd = exec.Command("gh", args...)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to create pull request %s\n", err)
	}

	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	url := urlRegex.FindString(stdout.String())
	copyToClipboard(url)
	fmt.Println(url)
}
