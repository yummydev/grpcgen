package codegen

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ClientGenerator generates code for connecting to a grpc client
// It is also the "base class" for other generators as they share its methods
type ClientGenerator struct {
	// Whether to Append or Create a file
	mode GeneratorMode
	// The target passed in
	variable *types.Var
	// The interface of the grpc client we are generating for
	inter *types.Interface
	// The pkg that the interface hails from
	interPkg *types.Package
}

// NewClientGenerator constructor
func NewClientGenerator(variable *types.Var, inter *types.Interface, interPkg *types.Package) *ClientGenerator {
	return &ClientGenerator{
		variable: variable,
		inter:    inter,
		interPkg: interPkg,
	}
}

// IsVarPkgDefined defined in the same package as the variable
func (g *ClientGenerator) IsVarPkgDefined(name string) error {
	pkg := g.variable.Pkg()
	obj := pkg.Scope().Lookup(name)
	if obj != nil {
		return fmt.Errorf("%q already exist", name)
	}
	return nil
}

// HandlerGenerator creates a handler generator for that particular
// method of the interface
func (g *ClientGenerator) HandlerGenerator(i int) *HandlerGenerator {
	return &HandlerGenerator{
		ClientGenerator: *g,
		MethodID:        i,
		Method:          g.inter.ExplicitMethod(i).Type().(*types.Signature),
	}
}

// HandlerGenerators creates a handler generator for every explicit method
func (g *ClientGenerator) HandlerGenerators() []*HandlerGenerator {
	gens := make([]*HandlerGenerator, g.inter.NumExplicitMethods())
	for i := range gens {
		gens[i] = g.HandlerGenerator(i)
	}
	return gens
}

// TypeName returns the name of type of the variable
// NOT the underlying type if aliased.
func (g *ClientGenerator) TypeName() string {
	return g.variable.Type().String()
}

// TypeNameStripped strips the package path from the type
func (g *ClientGenerator) TypeNameStripped() string {
	segments := strings.Split(g.TypeName(), ".")
	return segments[len(segments)-1]
}

// ID returns the variable Id as referenced by other packages
func (g *ClientGenerator) ID() string {
	return g.variable.Id()
}

// Name returns the variable name in the package scope
func (g *ClientGenerator) Name() string {
	return g.variable.Name()
}

// VarPkgName returns the package name of the variable
func (g *ClientGenerator) VarPkgName() string {
	return g.variable.Pkg().Name()
}

// VarPkgPath returns the package IMPORT path of the variable
func (g *ClientGenerator) VarPkgPath() string {
	return g.variable.Pkg().Path()
}

// InterPkgName returns the package name of the interface
func (g *ClientGenerator) InterPkgName() string {
	return g.interPkg.Name()
}

// InterPkgName returns the package IMPORT path of the interface
// removes vendor if found
func (g *ClientGenerator) InterPkgPath() string {
	p := g.interPkg.Path()
	stripped := strings.Split(p, "vendor/")
	return stripped[len(stripped)-1]
}

// PkgPath path to folder of the package relative to the GOPATH
func (g *ClientGenerator) PkgPath() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", g.VarPkgPath())
}

// GetFile attempts to open a file in append mode or create it.
// If noAppend is passed it will return EAlreadyExist.
// It will also set the mode of the generator, so not thread safe if using the
// same generator.
func (g *ClientGenerator) GetFile(fileName string, noAppend bool) (*os.File, error) {
	absFileName := filepath.Join(g.PkgPath(), fileName)
	var f *os.File
	var err error
	// if file doesn't exist create it
	if _, err = os.Stat(absFileName); os.IsNotExist(err) {
		f, err = os.Create(absFileName)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not create file %q", absFileName))
		}
		g.mode = GCreate
	} else {
		// don't attempt to open it if noAppend was passed
		if noAppend {
			return nil, EAlreadyExist
		}
		g.mode = GAppend
		f, err = os.OpenFile(absFileName, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Cannot open file %q for writing", absFileName))
		}
	}
	return f, nil
}

// GenerateClientsFile generates InitClient func if it doesn't alread exist
func (g *ClientGenerator) GenerateClientsFile() error {
	fileName := "clients.go"
	// check if the file doesn't exist
	f, err := g.GetFile(fileName, false)
	if err != nil {
		return err
	}
	defer f.Close()
	// if the file already exist make sure the funtion we are generating
	// doesn't already exist.
	if err := g.IsVarPkgDefined(g.InitClientFuncName()); err != nil {
		log.Println(err.Error() + " Skipping.")
		return nil
	}
	var s string
	if g.mode == GAppend {
		s, err = g.GenerateInitClientSnippet()
	} else {
		s, err = g.GenerateInitClient()
	}
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	return nil
}

// GenerateRunTests generates shell script to run the test
func (g *ClientGenerator) GenerateRunTestsFile() error {
	fileName := "run_test.sh"
	// check if the file doesn't exist
	f, err := g.GetFile(fileName, true)
	if err != nil {
		return err
	}
	defer f.Close()

	var s string
	s, err = g.GenerateRunTests()
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	return nil
}

// Available templates for generation

// GenerateRunTests generates tmpls/templates/run_test.tmpl
func (g *ClientGenerator) GenerateRunTests() (string, error) {
	return Generate(g, "run_test")
}

// GenerateInitClientSnippet generates tmpls/templates/init_client.tmpl
func (g *ClientGenerator) GenerateInitClientSnippet() (string, error) {
	return Generate(g, "init_client")
}

// GenerateInitClientSnippet generates tmpls/templates/init_client_full.tmpl
func (g *ClientGenerator) GenerateInitClient() (string, error) {
	return Generate(g, "init_client_full")
}
