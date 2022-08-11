package slackinviter

import (
	"os"
	"reflect"

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

// resolveEnv resolves any environment variables specified in the config.  it
// uses unholy method to iterate through struct fields and call os.ExpandEnv
// for each of them.
func (f *Fields) resolveEnv() {
	v := reflect.ValueOf(f)
	p := v.Elem()
	if p.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < p.NumField(); i++ {
		fld := p.Field(i)
		// will only resove strings
		if fld.Kind() != reflect.String {
			continue
		}

		// don't want any trouble, Sir.
		if !fld.CanSet() || !fld.IsValid() {
			continue
		}

		fld.SetString(os.ExpandEnv(fld.String()))
	}
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
	fld.resolveEnv() // resolve any environment variables
	return fld, nil
}
