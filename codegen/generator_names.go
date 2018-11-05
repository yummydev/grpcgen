package codegen

import (
	"fmt"
)

// InitClientFuncName name of function for setting up grpc connection
func (g *ClientGenerator) InitClientFuncName() string {
	return fmt.Sprintf("Init%s", g.TypeNameStripped())
}

// HandlerTestFuncName is the name of the test function to generate
func (g *HandlerGenerator) HandlerTestFuncName() string {
	return fmt.Sprintf("Test_%s_%s", g.TypeNameStripped(), g.HandlerName())
}
