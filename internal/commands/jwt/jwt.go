package jwt

import (
	"github.com/edsonjaramillo/crpytid/internal/commands/password"
	"github.com/edsonjaramillo/crpytid/internal/flags"
	"github.com/urfave/cli/v2"
)

// Main JWT secrets command
var JWTSecretsCommand = &cli.Command{
	Name:        "jwt",
	Usage:       "Generate JWT secrets",
	Subcommands: []*cli.Command{hs256subcommand, hs384subcommand, hs512subcommand},
}

// Subcommands
var hs256subcommand = &cli.Command{
	Name:    "hs256",
	Usage:   "Generate HS256 secrets",
	Flags:   []cli.Flag{flags.NoConsoleFlag, flags.NoClipboardFlag},
	Aliases: []string{"HS256", "256"},
	Action:  hs256Action,
}

var hs384subcommand = &cli.Command{
	Name:    "hs384",
	Usage:   "Generate HS384 secrets",
	Flags:   []cli.Flag{flags.NoConsoleFlag, flags.NoClipboardFlag},
	Aliases: []string{"HS384", "384"},
	Action:  hs384Action,
}

var hs512subcommand = &cli.Command{
	Name:    "hs512",
	Usage:   "Generate HS512 secrets",
	Flags:   []cli.Flag{flags.NoConsoleFlag, flags.NoClipboardFlag},
	Aliases: []string{"HS512", "512"},
	Action:  hs512Action,
}

// Actions
var hs256Action = func(cCtx *cli.Context) error {
	noConsole := cCtx.Bool("no-console")
	noClipboard := cCtx.Bool("no-clipboard")

	hs256Secret := password.GenerateRandom(32, false, false)

	flags.NoConsolePrinter(noConsole, hs256Secret)

	flags.ClipboardPrinter(noClipboard, hs256Secret)

	return nil
}

var hs384Action = func(cCtx *cli.Context) error {
	noConsole := cCtx.Bool("no-console")
	noClipboard := cCtx.Bool("no-clipboard")

	hs384Secret := password.GenerateRandom(48, false, false)

	flags.NoConsolePrinter(noConsole, hs384Secret)

	flags.ClipboardPrinter(noClipboard, hs384Secret)

	return nil
}

var hs512Action = func(cCtx *cli.Context) error {
	noConsole := cCtx.Bool("no-console")
	noClipboard := cCtx.Bool("no-clipboard")

	hs512Secret := password.GenerateRandom(64, false, false)

	flags.NoConsolePrinter(noConsole, hs512Secret)

	flags.ClipboardPrinter(noClipboard, hs512Secret)

	return nil
}
