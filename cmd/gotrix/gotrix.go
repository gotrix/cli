package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	appName string
	fs      *flag.FlagSet

	in       = &flags{}
	commands = []*command{
		{
			name: "create-app",
			exec: cmdCreateApp,
			desc: "Create a new application.",
		},
		{
			name: "create-component",
			exec: cmdCreateComponent,
			desc: "Create a new component.",
		},
		{
			name: "build-components",
			exec: cmdBuildComponents,
			desc: `Build components so files. Path can be modified by -path option. The defaut is "./components".`,
		},
	}
)

func main() {
	appName = "gotrix"
	fs = flag.NewFlagSet(appName, flag.ContinueOnError)
	fs.Usage = usage
	fs.BoolVar(&in.noColor, "no-color", false, "Do not colorize output (default false).")
	fs.BoolVar(&in.quiet, "quiet", false, "Do not print any output (default false).")
	fs.BoolVar(&in.help, "help", false, "Show this help.")
	fs.StringVar(&in.path, "path", "", "Path to contents directory (default depends on the command).")
	fs.Parse(os.Args[1:])
	cmdName := fs.Arg(0)
	if cmdName == "" || in.help == true {
		fs.Usage()
		return
	}
	cmd, err := getCommand(cmdName)
	if err != nil {
		echoErr(err)
	}
	if err := cmd.exec(); err != nil {
		echoErr(err)
	}
}

func cmdCreateApp() error {
	name := fs.Arg(1)
	if name == "" {
		fprintf(os.Stderr, "Usage:\n\n  %s %s [application name]\n",
			Yellow.SPrint(appName),
			Blue.SPrint("create-app"))
		os.Exit(2)
	}
	echo(Blue, `Creating new application "%s".`, name)
	path, err := getPath(".")
	if err != nil {
		return err
	}
	checkoutPath := filepath.Join(path, name)
	repo := "github.com/gotrix/skull"
	cmd := exec.Command("git", "clone",
		"https://"+repo, checkoutPath)
	outBuf := bytes.NewBuffer([]byte{})
	cmd.Stderr = outBuf
	cmd.Stdout = outBuf
	echo(Blue, `Cloning %s to %s`, repo, checkoutPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(`failed to clone %s: %v: %s`,
			repo, err, strings.TrimRight(outBuf.String(), "\n"))
	}
	echo(Green, `Cloned %s to %s`, repo, checkoutPath)
	outBuf.Reset()
	cmd = exec.Command("git", "remote", "remove", "origin")
	cmd.Stderr = outBuf
	cmd.Stdout = outBuf
	cmd.Dir = checkoutPath
	echo(Blue, "Removing origins")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf(`failed remove origins: %v: %s`,
			err, strings.TrimRight(outBuf.String(), "\n"))
	}
	echo(Green, "Removed origins")
	return nil
}

func cmdCreateComponent() error {
	name := fs.Arg(1)
	if name == "" {
		fprintf(os.Stderr, "Usage:\n\n  %s %s [component name]\n",
			Yellow.SPrint(appName),
			Blue.SPrint("create-app"))
		os.Exit(2)
	}
	echo(Blue, `Creating new component "%s".`, name)
	return nil
}

func cmdBuildComponents() error {
	path, err := getPath("./components")
	if err != nil {
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
				name, err, strings.TrimRight(outBuf.String(), "\n"))
		}
		echo(Green, `successfully built "%s" to %s`, name, out)
	}
	return nil
}

func getPath(def string) (string, error) {
	var (
		path = in.path
		err  error
	)
	if path == "" {
		path = def
	}
	if strings.HasPrefix(path, "~") {
		hd, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = hd + path[1:]
	}
	path, err = filepath.Abs(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf(`%s: no such file or directory`, path)
		}
		return "", err
	}
	return path, nil
}

func getCommand(name string) (*command, error) {
	for _, cmd := range commands {
		if cmd.name == name {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf(`unknown command "%s"`, name)
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

func fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func echo(dot color, format string, args ...interface{}) {
	if !in.quiet {
		fprintf(os.Stdout, fmt.Sprintf("%s ", dot.SPrint("*"))+"%s\n", fmt.Sprintf(format, args...))
	}
}

func echoErr(err error) {
	if !in.quiet {
		fprintf(os.Stdout, fmt.Sprintf("%s ", Red.SPrint("* Error: "))+"%s\n", err.Error())
	}
	os.Exit(1)
}

func multiSuffix(l int, s string) string {
	if l == 1 {
		return s
	}
	return s + "s"
}

func makeDesc(s string) string {
	parts := strings.Split(strings.Trim(s, " "), " ")
	res := make([]string, int(math.Ceil(float64(len(s)/60)))+1)
	line := 0
	for _, p := range parts {
		if len(res[line]+" "+p) <= 60 {
			if res[line] == "" {
				res[line] += p
			} else {
				res[line] += " " + p
			}
		} else {
			line++
			res[line] += p
		}
	}
	ret := make([]string, 0, len(res))
	for _, v := range res {
		if v != "" {
			ret = append(ret, v)
		}
	}
	return strings.Join(ret, "\n        ")
}

func usage() {
	fprintf(os.Stdout, "Usage:\n")
	fprintf(os.Stdout, "\n  %s [%s] %s\n", Yellow.SPrint(appName), Green.SPrint("options"), Blue.SPrint("command"))
	fprintf(os.Stdout, "\nCommands:\n\n")
	step := "        "

	for _, cmd := range commands {
		fprintf(os.Stdout, "  %s\n%s%s\n\n", Blue.SPrint(cmd.name), step, makeDesc(cmd.desc))
	}
	fprintf(os.Stdout, "Options:\n")
	fs.VisitAll(func(f *flag.Flag) {
		_, usage := flag.UnquoteUsage(f)
		fprintf(os.Stdout, "\n  %s\n", Green.SPrint("-"+f.Name))
		fprintf(os.Stdout, "        %s\n", makeDesc(usage))
	})
}

type flags struct {
	noColor bool
	quiet   bool
	help    bool
	path    string
}

type command struct {
	name string
	desc string
	exec func() error
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
	// Black color.
	Black color = "\u001b[30m"

	// Red color.
	Red color = "\u001b[31m"

	// Green color.
	Green color = "\u001b[32m"

	// Yellow color.
	Yellow color = "\u001b[33m"

	// Blue color.
	Blue color = "\u001b[34m"

	// Reset color.
	Reset color = "\u001b[0m"
)
