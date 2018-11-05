package codegen

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

// HandlerGenerator represents a specific handler of the
// ClientGenerator interface
type HandlerGenerator struct {
	MethodID int
	Method   *types.Signature
	ClientGenerator
}

// HandlerName returns the name of the given handler
func (g *HandlerGenerator) HandlerName() string {
	return g.ClientGenerator.inter.ExplicitMethod(g.MethodID).Name()
}

// ReqName returns the name of the type of the Request paramater
func (g *HandlerGenerator) ReqName() string {
	return g.Method.Params().At(1).Type().String()
}

// ReqNameStripped returns the name of the type of the Request paramater - package name
func (g *HandlerGenerator) ReqNameStripped() string {
	segments := strings.Split(g.Method.Params().At(1).Type().String(), ".")
	return segments[len(segments)-1]
}

// ResName returns the name of the type of the Response parameter
func (g *HandlerGenerator) ResName() string {
	return g.Method.Results().At(0).Type().String()
}

// ResNameStripped returns the name of the type of the Response parameter - package name
func (g *HandlerGenerator) ResNameStripped() string {
	segments := strings.Split(g.Method.Results().At(0).Type().String(), ".")
	return segments[len(segments)-1]
}

// GenerateHandlerTestFile fully generates the handler test file
// Since test aren't included with the package scope, we are unable
// to detect if they already exist, we will not append here
func (g *HandlerGenerator) GenerateHandlerTestFile() error {
	fileName := strcase.ToSnake(fmt.Sprintf("%s_%s_test", g.TypeNameStripped(), g.HandlerName())) + ".go"
	// check if the file doesn't exist
	f, err := g.GetFile(fileName, true)
	if err != nil {
		if err == EAlreadyExist {
			return nil
		}
		return err
	}
	defer f.Close()
	// if the file already exist make sure the funtion we are generating
	// doesn't already exist.
	if err := g.IsVarPkgDefined(g.HandlerTestFuncName()); err != nil {
		log.Println(err.Error() + " Skipping.")
		return nil
	}

	s, err := g.GenerateHandlerTest()
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	return nil
}

// GenerateHandlerTestSnippet generates just the test function
func (g *HandlerGenerator) GenerateHandlerTestSnippet() (string, error) {
	return Generate(g, "handler_test")
}

// GenerateHandlerTest generates the entire file
func (g *HandlerGenerator) GenerateHandlerTest() (string, error) {
	return Generate(g, "handler_test_full")
}
