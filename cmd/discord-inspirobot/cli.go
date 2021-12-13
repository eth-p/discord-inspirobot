package main

import (
	"flag"

	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

var Cli = cli.App{
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "Klog verbosity level.",
			DefaultText: "2",
		},
		&cli.StringFlag{
			Name:     "bot-token",
			Usage:    "The Discord bot token.",
			EnvVars:  []string{"DISCORD_TOKEN"},
			FilePath: ".discord_token",
			Required: true,
		},
	},

	Before: func(context *cli.Context) error {
		// Set up klog verbosity flags.
		klogOptions := flag.FlagSet{}
		klog.InitFlags(&klogOptions)
		return klogOptions.Set("v", context.String("verbose"))
	},
}
