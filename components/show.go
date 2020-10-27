package components

import (
	"github.com/mkideal/cli"
)

// ShowT struct.
type ShowT struct {
	cli.Helper2
	Verbose bool `cli:"verbose" usage:"More output"`
}

// Show is showing a list of components.
func Show(ctx *cli.Context) error {
	argv := ctx.Argv().(*ShowT)
	path, err := getPath()
	if err != nil {
		return err
	}
	l, err := List(path)
	if err != nil {
		return err
	}
	for _, c := range l {
		if argv.Verbose {
			ctx.String("\nComponent \"%s\"\n", ctx.Color().Green(c.Name))
			ctx.String("  %s\n", ctx.Color().Blue("JavaScript:"))
			ctx.String("    %s %d files\n", ctx.Color().Bold("Async:"), len(c.AsyncJsPaths))
			for _, s := range c.AsyncJsPaths {
				ctx.String("      %s\n", s)
			}
			ctx.String("    %s %d files\n", ctx.Color().Bold("Defer:"), len(c.DeferJsPaths))
			for _, s := range c.DeferJsPaths {
				ctx.String("      %s\n", s)
			}
			ctx.String("    %s %d files\n", ctx.Color().Bold("Regular:"), len(c.JsPaths))
			for _, s := range c.JsPaths {
				ctx.String("      %s\n", s)
			}
			ctx.String("  %s\n", ctx.Color().Blue("Stylesheets:"))
			ctx.String("    %s %d files\n", ctx.Color().Bold("Async:"), len(c.AsyncCssPaths))
			for _, s := range c.AsyncCssPaths {
				ctx.String("      %s\n", s)
			}
			ctx.String("    %s %d files\n", ctx.Color().Bold("Regular:"), len(c.CssPaths))
			for _, s := range c.CssPaths {
				ctx.String("      %s\n", s)
			}
		} else {
			ctx.String("%s\n", c.Name)
		}
	}
	ctx.String("\n")
	return nil
}
