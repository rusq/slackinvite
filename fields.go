package slackinviter

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Fields struct {
	//body
	SlackWorkspace string `yaml:"slack_workspace,omitempty"`
	//form
	SubmitButton string `yaml:"submit_button,omitempty"`
	// footer
	Website      string `yaml:"website,omitempty"`
	Copyright    string `yaml:"copyright,omitempty"`
	TelegramLink string `yaml:"telegram_link,omitempty"`
	GithubLink   string `yaml:"github_link,omitempty"`
}

func LoadFields(filename string) (Fields, error) {
	var fld Fields
	f, err := os.Open(filename)
	if err != nil {
		return fld, err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f, yaml.DisallowUnknownField())
	if err := dec.Decode(&fld); err != nil {
		return fld, err
	}
	return fld, nil
}
