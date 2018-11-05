package tmpls

import (
	"text/template"

	"github.com/gobuffalo/packr"
)

// Templates contains all the package templates
var Templates packr.Box

func init() {
	Templates = packr.NewBox("./templates")
}

func GetTemplate(baseName string) (*template.Template, error) {
	tString, err := Templates.FindString(baseName + ".tmpl")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(baseName).Parse(tString)
	return tmpl, err
}
