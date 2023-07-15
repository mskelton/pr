package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func branchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getDefaultTitle() string {
	branch, err := branchName()
	if err != nil {
		return ""
	}

	branch = strings.ReplaceAll(branch, "-", " ")

	re := regexp.MustCompile(`^fcs `)
	branch = re.ReplaceAllString(branch, "FCS-")

	parts := strings.SplitN(branch, " ", 2)

	if len(parts) > 1 && len(parts[1]) > 0 {
		parts[1] = strings.ToUpper(parts[1][:1]) + parts[1][1:]
	}

	return strings.Join(parts, " ")
}

func buildSurvey() []*survey.Question {
	return []*survey.Question{
		{
			Name:     "title",
			Prompt:   &survey.Input{Message: "Title:", Default: getDefaultTitle()},
			Validate: survey.Required,
		},
		{
			Name:   "body",
			Prompt: &survey.Input{Message: "Body:"},
		},
	}
}

func main() {
	qs := buildSurvey()
	answers := struct {
		Title string
		Body  string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
