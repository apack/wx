package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/apack/wx/internal/compiler"
)

//go:embed gen.gotmpl
var genTmplStr string

var genTmpl = template.Must(template.New("gen").Funcs(template.FuncMap{
	"printBinary": printBinary,
}).Parse(genTmplStr))

type Generator struct {
	views []*compiler.CompiledView
	dir   string
}

type GenerateConfig struct {
	Package string
	Views   []*compiler.CompiledView
}

func NewGenerator(views []*compiler.CompiledView, dir string) *Generator {
	return &Generator{views: views, dir: dir}
}

func (g *Generator) Generate(w io.Writer) error {
	dir, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}
	dir = filepath.Base(dir)
	return genTmpl.Execute(w, GenerateConfig{
		Package: strings.ToLower(dir),
		Views:   g.views,
	})
}

func printBinary(b []byte) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[]byte{")
	for i := 0; i < len(b); i++ {
		if i != 0 {
			buf.WriteString(", ")
		}
		if i%12 == 0 {
			buf.WriteString("\n\t")
		}
		buf.WriteString(fmt.Sprintf("0x%02X", int(b[i])))
	}
	if len(b) > 0 {
		buf.WriteString(",\n")
	}
	buf.WriteString("}")
	return buf.String()
}
