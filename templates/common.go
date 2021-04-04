package templates

import (
	"io"
	"text/template"
)

func Compile(tpl string, out io.Writer, data interface{}) error {
	// Supposing that  The template name does not matter
	genTemp := template.New("Template")
	ctemplate, err := genTemp.Parse(tpl)
	if err != nil {
		return err
	}
	if err = ctemplate.Execute(out, data); err != nil {
		return err
	}
	return nil
}
