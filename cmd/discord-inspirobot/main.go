package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/eth-p/discord-inspirobot/internal/bot"
)

var Commands []bot.CommandRegistration

func init() {
	Cli.Commands = append(Cli.Commands, &cli.Command{
		Name:  "run",
		Usage: "Run the bot.",
		Action: func(context *cli.Context) error {
			instance, err := bot.New(context.String("bot-token"), context.Context, bot.BotCommands(Commands...))
			if err != nil {
				return err
			}

			return instance.Wait()
		},
	})
}

func main() {
	err := Cli.RunContext(sigintContext(), os.Args)
	if err != nil {
		println(err.Error())
	}
}
