package codegen

import (
	"bytes"

	"github.com/yummydev/grpcgen/tmpls"
)

// GeneratorMode determines the way file generation will work
type GeneratorMode = string

const (
	// GCreate create file
	GCreate GeneratorMode = "CREATE"
	// GAppend append to end of file
	GAppend GeneratorMode = "APPEND"
)

// Generate retrieves the template and executes it with the generator
func Generate(generator interface{}, handlerName string) (string, error) {
	tmpl, err := tmpls.GetTemplate(handlerName)
	if err != nil {
		return "", err
	}
	b := &bytes.Buffer{}
	err = tmpl.Execute(b, generator)
	if err != nil {
		return "", err
	}
	return string(b.Bytes()), nil
}
