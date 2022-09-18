package ssr

import (
	"encoding/json"
	"fmt"

	_ "embed"

	"github.com/apack/wx/internal/js"
)

func Load(vm js.VM, code string) (*SSR, error) {
	if err := vm.Script("ssr.js", code); err != nil {
		return nil, err
	}
	// TODO make dev configurable
	return &SSR{vm, true}, nil
}

type SSR struct {
	VM  js.VM
	Dev bool
}

func (c *SSR) Execute(dst interface{}, props interface{}) error {
	propsData, err := json.Marshal(props)
	if err != nil {
		return err
	}
	expr := fmt.Sprintf(`;JSON.stringify(__ssr__.default.render(%s))`, string(propsData))
	result, err := c.VM.Eval("ssr.js", expr)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(result), dst); err != nil {
		return err
	}
	return nil
}
