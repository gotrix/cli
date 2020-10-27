package components

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	gopath "path"
	"path/filepath"
	"strings"

	"github.com/gotrix/cli/helpers"
	"github.com/gotrix/cli/mod"
	"github.com/mkideal/cli"
	"github.com/tdewolff/minify/v2"
	minCss "github.com/tdewolff/minify/v2/css"
	minHtml "github.com/tdewolff/minify/v2/html"
	minJs "github.com/tdewolff/minify/v2/js"
)

// BuildT struct.
type BuildT struct {
	cli.Helper2
	Name []string `cli:"name" usage:"Component name"`
}

// Build is building one or more components.
func Build(ctx *cli.Context) error {
	argv := ctx.Argv().(*BuildT)
	path, err := getPath()
	if err != nil {
		return err
	}
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	if err := helpers.Rm(filepath.Join(root, ".gotrix")); err != nil {
		return err
	}
	if len(argv.Name) > 0 {
		for _, name := range argv.Name {
			info, err := getInfo(name, filepath.Join(path, name))
			if err != nil {
				return err
			}
			if err := build(ctx, info); err != nil {
				return err
			}
		}
	} else {
		list, err := List(path)
		if err != nil {
			return err
		}
		for _, info := range list {
			if err := build(ctx, info); err != nil {
				return err
			}
		}
	}
	ctx.String("Done\n")
	return nil
}

func build(ctx *cli.Context, info *Info) error {
	ctx.String("Building component %s\n", info.Name)
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	hash := getHash()
	staticPath := filepath.Join(root, ".gotrix", "static", "component", hash)
	if err := helpers.MakeDir(staticPath); err != nil {
		return err
	}
	m := minify.New()
	m.Add("text/html", &minHtml.Minifier{
		KeepConditionalComments: true,
		KeepDefaultAttrVals:     true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
		KeepWhitespace:          false,
	})
	m.AddFunc("text/css", minCss.Minify)
	m.AddFunc("application/javascript", minJs.Minify)
	copyFile := func(f, mt string) error {
		src := f
		dst := filepath.Join(staticPath, filepath.Base(f))
		sourceFileStat, err := os.Stat(src)
		if err != nil {
			return err
		}
		if !sourceFileStat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", src)
		}
		source, err := os.Open(src)
		if err != nil {
			return err
		}
		defer func() { _ = source.Close() }()
		destination, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer func() { _ = destination.Close() }()
		if err := m.Minify(mt, destination, source); err != nil {
			return err
		}
		return nil
	}
	toURL := func(f string) string {
		return gopath.Join("/gotrix/component", hash, f)
	}
	t, err := ioutil.ReadFile(info.TemplatePath)
	if err != nil {
		return err
	}
	t, err = m.Bytes("text/html", t)
	if err != nil {
		return err
	}
	pack, err := mod.FromGoMod("go.mod")
	if err != nil {
		return err
	}
	td := &wrapTemplateData{
		Hash:     hash,
		Package:  pack + "/components/" + info.Name,
		Template: string(t),
		AsyncCSS: filesList{},
		CSS:      filesList{},
		AsyncJS:  filesList{},
		DeferJS:  filesList{},
		JS:       filesList{},
	}
	for _, s := range info.Statics() {
		base := filepath.Base(s)
		if strings.HasSuffix(s, ".async.js") {
			if err := copyFile(s, "application/javascript"); err != nil {
				return err
			}
			td.AsyncJS = append(td.AsyncJS, toURL(base))
			continue
		}
		if strings.HasSuffix(s, ".defer.js") {
			if err := copyFile(s, "application/javascript"); err != nil {
				return err
			}
			td.DeferJS = append(td.DeferJS, toURL(base))
			continue
		}
		if strings.HasSuffix(s, ".js") {
			if err := copyFile(s, "application/javascript"); err != nil {
				return err
			}
			td.JS = append(td.JS, toURL(base))
			continue
		}
		if strings.HasSuffix(s, ".async.css") {
			if err := copyFile(s, "text/css"); err != nil {
				return err
			}
			td.AsyncCSS = append(td.AsyncCSS, toURL(base))
			continue
		}
		if strings.HasSuffix(s, ".css") {
			if err := copyFile(s, "text/css"); err != nil {
				return err
			}
			td.CSS = append(td.CSS, toURL(base))
			continue
		}
	}
	res := bytes.NewBuffer([]byte{})
	if err := wrapTemplate.Execute(res, td); err != nil {
		return err
	}
	buildPath := filepath.Join(root, ".gotrix", "build", hash)
	if err := helpers.MakeDir(buildPath); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(buildPath, "main.go"), res.Bytes(), 0644); err != nil {
		return err
	}
	outBuf := bytes.NewBuffer([]byte{})
	cmd := exec.Command("go", "build",
		"-buildmode=plugin",
		"-o", filepath.Join(info.Dir, "component.so"),
		filepath.Join(buildPath, "main.go"))
	cmd.Stdout = outBuf
	cmd.Stderr = outBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(`failed to build component%s: %v: %s`,
			info.Name, err, strings.TrimRight(outBuf.String(), "\n"))
	}
	if err := helpers.Rm(buildPath); err != nil {
		return err
	}
	return nil
}
