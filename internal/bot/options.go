package bot

// Option is an option for the Bot.
type Option func(bot uninitializedBot) error

// BotCommands adds slash commands to the Bot.
func BotCommands(commands ...CommandRegistration) Option {
	return func(bot uninitializedBot) error {
		for _, command := range commands {
			bot.commands[command.spec().Name] = command
		}

		return nil
	}
}
