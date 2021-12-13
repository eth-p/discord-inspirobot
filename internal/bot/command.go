package bot

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is the function used to handle a slash command.
type CommandHandler func(ctx *CommandContext, interaction Interaction)

// CommandRegistration is an interface providing the required data to register a Discord slash command.
type CommandRegistration interface {
	spec() *discordgo.ApplicationCommand
	handler() CommandHandler
	componentHandlers() map[string]ComponentHandler
}

// ---------------------------------------------------------------------------------------------------------------------

// Command is a struct for registering a Bot slash command.
// This implements CommandRegistration, and can either be used as-is, or composed inside another struct.
//
// Example:
//   bot.New("MY_TOKEN", context.Background(), bot.BotCommands(&bot.Command {
//     Spec: &discordgo.ApplicationCommand {
//       Name: "ping",
//       Description: "replies with pong",
//     },
//     Handler: func(ctx *bot.CommandContext, i bot.Interaction) error {
//       i.Respond(&discordgo.InteractionResponseData {
//         Content: "Hey there! Nice to meet you!",
//       }
//     }
//   })
type Command struct {
	Spec              *discordgo.ApplicationCommand
	Handler           CommandHandler
	ComponentHandlers map[string]ComponentHandler
}

func (c *Command) spec() *discordgo.ApplicationCommand {
	return c.Spec
}

func (c *Command) handler() CommandHandler {
	return c.Handler
}

func (c *Command) componentHandlers() map[string]ComponentHandler {
	return c.ComponentHandlers
}

// Name returns the command name.
func (c *Command) Name() string {
	return c.Spec.Name
}

// ---------------------------------------------------------------------------------------------------------------------

// CommandContext is a context.Context passed to a CommandHandler.
// This contains additional data about the bot and the command instance.
type CommandContext struct {
	exceptional
	botContext

	// Self is the pointer to the "Command" struct registered.
	// This may be used with type assertions to get the original struct back.
	//
	// Example:
	//   func MyHandler(ctx *CommandContext, interaction bot.Interaction) error {
	//     self := ctx.Self.(*MyCommand)
	//     self.AnyCustomField = 3
	//   }
	Self interface{}
}

// Command returns a Command struct that represents the command being executed.
func (c *CommandContext) Command() Command {
	if cmd, ok := c.Self.(*Command); ok {
		return *cmd
	}

	// Synthesize a Command instance.
	cmdIface := c.Self.(CommandRegistration)
	return Command{
		Spec:    cmdIface.spec(),
		Handler: cmdIface.handler(),
	}
}
