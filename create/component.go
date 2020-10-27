package create

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gotrix/cli/helpers"
	"github.com/mkideal/cli"
)

// ComponentT struct.
type ComponentT struct {
	Name string `cli:"*name" usage:"Component name"`
}

// Component creates a new component.
func Component(ctx *cli.Context) error {
	argv := ctx.Argv().(*ComponentT)
	ctx.String("Creating %s component\n",
		ctx.Color().Green(argv.Name))
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	path := filepath.Join(root, "components", argv.Name)
	if err := helpers.MakeDir(path); err != nil {
		return err
	}
	for f, c := range componentFiles {
		ctx.String("Creating %s file\n",
			ctx.Color().Blue(f))
		v := strings.ReplaceAll(c, "%NAME%", argv.Name)
		if err := ioutil.WriteFile(filepath.Join(path, f), []byte(v), 0644); err != nil {
			return err
		}
	}
	ctx.String("Done\n")
	return nil
}

var (
	componentFiles = map[string]string{
		"component.go":        componentGo,
		"component.html":      componentTemplate,
		"component.async.css": componentCss,
		"component.async.js":  componentJs,
	}
)

const (
	componentGo = `package component

import (
	"github.com/gotrix/gotrix"
)

type Component struct{}

// Component implements gotrix.Component.
func (*Component) Component(ctx gotrix.Context, params gotrix.ComponentParams) (map[string]interface{}, error) {

	// Do something logical...

	return map[string]interface{}{
		"name": "Simon",
	}, nil
}`
	componentTemplate = `<div class="%NAME%-component">
	<h2>Hey there</h2>
	<p>Here's your new component</p>
	<p>Now, go do something special...</p>
	<p>Kind Regards,<br/>{{.Data("name")}}</p>
</div>`

	componentCss = `.%NAME%-component {
	font-family: sans-serif;
	color: #333;
}`

	componentJs = `(function() {
	console.log('component %NAME% script loaded');
})();`
)
