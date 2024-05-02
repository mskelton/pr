package main

import (
	"fmt"
	"log"
	"os"
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
	flags := os.Args[1:]

	branch := branchName()
	prefixes := getPrefixes()
	prefix, title := getDefaultTitle(prefixes, branch)
	qs := buildSurvey(title)

	answers := struct {
		Title string
		Body  string
	}{}

	err := survey.Ask(qs, &answers)
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
	fmt.Println("gh pr create", strings.Join(flags, " "), "--title", quote(finalTitle), "--body", quote(answers.Body))
	fmt.Println()

	// Push first to ensure the command will succeed
	cmd := exec.Command("git", "push")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to push branch %s\n", err)
	}

	args := []string{"pr", "create", "--title", finalTitle, "--body", answers.Body}
	args = append(args, flags...)

	cmd = exec.Command("gh", args...)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to create pull request %s\n", err)
	}

	fmt.Println("Pull request created!")
}
