package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	flag.BoolVar(&in.noColor, "no-color", false, "do not colorize output")
	flag.BoolVar(&in.quiet, "quiet", false, "do not print any output")
	flag.StringVar(&in.path, "path", "", "path to use")
	flag.Parse()
	var (
		command = flag.Arg(0)
		err     error
	)
	if command == "" {
		flag.Usage()
		return
	}
	switch command {
	case "build-components":
		err = cmdBuildComponents()
	default:
		err = fmt.Errorf(`unkown command "%s"`, command)
	}
	if err != nil {
		echoErr(err)
	}
}

func cmdBuildComponents() error {
	var (
		path = in.path
		err  error
	)
	if path == "" {
		path = "./components"
	}
	path, err = filepath.Abs(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(`%s: no such file or directory`, path)
		}
		return err
	}
	files, err := readPathsRecursive(path, ".go")
	if err != nil {
		return fmt.Errorf("failed to read %s", path)
	}
	lf := len(files)
	if lf == 0 {
		return fmt.Errorf("no components found in %s", path)
	}
	getName := func(f string) string {
		parts := strings.Split(strings.TrimPrefix(f, path), "/")
		return strings.Trim(strings.Join(parts[:len(parts)-1], "/"), "/")
	}
	echo(Blue, "building components from %s", path)
	echo(Blue, "found %d %s", lf, multiSuffix(lf, "component"))
	for _, f := range files {
		name := Green.SPrint(getName(f))
		out := strings.TrimSuffix(f, ".go") + ".so"
		echo(Blue, `building component "%s" from %s`, name, f)
		outBuf := bytes.NewBuffer([]byte{})
		cmd := exec.Command("go", "build",
			"-buildmode=plugin",
			"-o", out, f)
		cmd.Stdout = outBuf
		cmd.Stderr = outBuf
		if err := cmd.Run(); err != nil {
			return fmt.Errorf(`failed to build "%s": %v: %s`,
				name, err, strings.TrimRight(outBuf.String(),"\n"))
		}
		echo(Green, `successfully built "%s" to %s`, name, out)
	}
	return nil
}

func readPathsRecursive(dir string, suffix string) (list []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			list = append(list, path)
		}
		return nil
	})
	return
}

func echo(dot color, format string, args ...interface{}) {
	if !in.quiet {
		_, _ = fmt.Fprint(os.Stdout,
			dot.SPrint("*"), " ", fmt.Sprintf(format, args...), "\n")
	}
}

func echoErr(err error) {
	if !in.quiet {
		_, _ = fmt.Fprint(os.Stderr,
			Red.SPrint("* error:"), " ", err.Error(), "\n")
	}
	os.Exit(1)
}

func multiSuffix(l int, s string) string {
	if l == 1 {
		return s
	}
	return s + "s"
}

var (
	in = &flags{}
)

type flags struct {
	noColor bool
	quiet   bool
	path    string
}

type color string

func (c color) String() string {
	return string(c)
}

func (c color) SPrint(s string) string {
	if in.noColor {
		return s
	}
	return fmt.Sprint(c, s, Reset)
}

const (
	Black  color = "\u001b[30m"
	Red    color = "\u001b[31m"
	Green  color = "\u001b[32m"
	Yellow color = "\u001b[33m"
	Blue   color = "\u001b[34m"
	Reset  color = "\u001b[0m"
)
