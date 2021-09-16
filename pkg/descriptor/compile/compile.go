package compile

import (
	"github.com/toitware/can-go/internal/generate"
)

type Result generate.CompileResult

func Run(sourceFile string, data []byte) (result *Result, err error) {
	res, err := generate.Compile(sourceFile, data)
	return (*Result)(res), err
}
