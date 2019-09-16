package client

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// User naive user data structure
type User struct {
	Name string
	ID   string `yaml:"id,omitempty"`
}

// Template templates a user struct onto a string
func (u *User) Template(t string) (string, error) {
	tmpl := template.Must(template.New("").Parse(t))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, u); err != nil {
		return "", errors.Wrap(err, "tmpl execute")
	}
	return b.String(), nil
}
