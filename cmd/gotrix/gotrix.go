package main

import (
	"fmt"
	"os"

	"github.com/gotrix/cli/components"
	"github.com/gotrix/cli/create"
	"github.com/mkideal/cli"
)

var (
	help    = cli.HelpCommand("Display help information")
	cmdRoot = &cli.Command{
		Desc: "GOTRIX CLI Tool",
		Fn: func(ctx *cli.Context) error {
			ctx.String("GOTRIX CLI Tool.\nRun 'gotrix help' for help.\n")
			return nil
		},
	}
	cmdCreate = &cli.Command{
		Name: "create",
		Desc: "Create new something",
		Fn: func(ctx *cli.Context) error {
			return help.Fn(ctx)
		},
	}
	cmdCreateApp = &cli.Command{
		Name: "app",
		Desc: "Create new application",
		Argv: func() interface{} { return new(create.AppT) },
		Fn:   create.App,
	}
	cmdCreateComponent = &cli.Command{
		Name: "component",
		Desc: "Create new component",
		Argv: func() interface{} { return new(create.ComponentT) },
		Fn:   create.Component,
	}
	cmdComponents = &cli.Command{
		Name: "components",
		Desc: "Manipulate Components",
		Fn: func(ctx *cli.Context) error {
			return help.Fn(ctx)
		},
	}
	cmdComponentsShow = &cli.Command{
		Name: "show",
		Desc: "Show Components",
		Argv: func() interface{} { return new(components.ShowT) },
		Fn:   components.Show,
	}
	cmdComponentsBuild = &cli.Command{
		Name: "build",
		Desc: "Build one or more components",
		Argv: func() interface{} { return new(components.BuildT) },
		Fn:   components.Build,
	}
)

func main() {
	if err := cli.Root(cmdRoot,
		cli.Tree(help),
		cli.Tree(cmdCreate,
			cli.Tree(cmdCreateApp),
			cli.Tree(cmdCreateComponent),
		),
		cli.Tree(cmdComponents,
			cli.Tree(cmdComponentsShow),
			cli.Tree(cmdComponentsBuild),
		),
	).Run(os.Args[1:]); err != nil {
		if _, e := fmt.Fprintln(os.Stderr, err); e != nil {
			panic(fmt.Sprintf("%v: %v", e, err))
		}
		os.Exit(1)
	}
}
