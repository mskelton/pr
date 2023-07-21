package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func branchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func getDefaultTitle(branch string) (string, string) {
	re := regexp.MustCompile(`^fcs-\d+`)
	prefix := re.FindString(branch)
	title := strings.TrimSpace(strings.ReplaceAll(re.ReplaceAllString(branch, ""), "-", " "))

	if len(title) > 0 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}

	return strings.ToUpper(prefix), title
}

func buildSurvey(title string) []*survey.Question {
	return []*survey.Question{
		{
			Name:     "title",
			Prompt:   &survey.Input{Message: "Title:", Default: title},
			Validate: survey.Required,
		},
		{
			Name:   "body",
			Prompt: &survey.Input{Message: "Body:"},
		},
	}
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func main() {
	branch := branchName()
	prefix, title := getDefaultTitle(branch)
	qs := buildSurvey(title)

	answers := struct {
		Title string
		Body  string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatal(err.Error())
	}

	finalTitle := answers.Title

	// Automatically add the ticket to the title if present
	if prefix != "" {
		finalTitle = prefix + " " + finalTitle
	}

	fmt.Println("Creating pull request...")
	fmt.Println("gh pr create --title", quote(finalTitle), "--body", quote(answers.Body))
	fmt.Println()

	// Push first to ensure the command will succeed
	cmd := exec.Command("git", "push")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to push branch %s\n", err)
	}

	cmd = exec.Command("gh", "pr", "create", "--title", finalTitle, "--body", answers.Body)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to create pull request %s\n", err)
	}

	fmt.Println("Pull request created!")
}
