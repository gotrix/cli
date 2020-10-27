package create

import (
	"fmt"
	"path/filepath"

	"github.com/gotrix/cli/helpers"
	"github.com/mkideal/cli"
)

// AppT struct.
type AppT struct {
	Name   string `cli:"*name" usage:"Application name"`
	Module string `cli:"module" usage:"Go module, defaults to --name"`
}

// App creates a new application.
func App(ctx *cli.Context) error {
	argv := ctx.Argv().(*AppT)
	if argv.Module == "" {
		argv.Module = argv.Name
	}
	ctx.String("Creating application %s (module: %s)\n",
		ctx.Color().Green(argv.Name),
		ctx.Color().Italic(argv.Module))
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	repo := "https://github.com/gotrix/skull"
	checkoutPath := filepath.Join(root, argv.Name)
	ctx.String("Cloning %s into %s\n",
		ctx.Color().Underline(repo),
		ctx.Color().Underline(checkoutPath))
	if err := helpers.Run("git", []string{"clone", repo, checkoutPath}, ""); err != nil {
		return fmt.Errorf("failed to create application: %v", err)
	}
	ctx.String("Cleaning remotes\n")
	if err = helpers.Run("git", []string{"remote", "remove", "origin"}, checkoutPath); err != nil {
		return fmt.Errorf(`failed to clean origins: %v`, err)
	}
	ctx.String("Updating module name\n")
	if err = helpers.Run("go", []string{"mod", "edit", "-module", argv.Module}, checkoutPath); err != nil {
		return fmt.Errorf(`failed to update module name: %v`, err)
	}
	ctx.String("Updating gotrix to latest version\n")
	if err = helpers.Run("go", []string{"get", "-u", "github.com/gotrix/gotrix"}, checkoutPath); err != nil {
		return fmt.Errorf(`failed to update module name: %v`, err)
	}
	ctx.String("Done\n")
	return nil
}
