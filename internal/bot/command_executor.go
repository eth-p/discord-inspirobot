package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// executeCommand executes a Command and returns its CommandContext and ExecutionStatus.
func executeCommand(ctx botContext, command CommandRegistration, interaction *discordgo.Interaction) (*CommandContext, ExecutionStatus) {
	cmdCtx := &CommandContext{
		botContext: ctx,
		Self:       command,
	}

	// Run the command handler.
	status := ExecutionStatus{}
	status.TimeStarted = time.Now()
	status.Err = cmdCtx.try(func() {
		command.handler()(cmdCtx, Interaction{
			context:     &cmdCtx.botContext,
			Interaction: interaction,
		})
	})
	status.TimeFinished = time.Now()

	// Return the context and status.
	return cmdCtx, status
}

// commandNotFound is an implementation of a CommandRegistration that returns an error.
// This is used to simplify command handling control flow.
type commandNotFound struct {
	name string
}

func (c commandNotFound) spec() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: c.name,
	}
}

func (c commandNotFound) handler() CommandHandler {
	return func(ctx *CommandContext, interaction Interaction) {
		ctx.Fatal(fmt.Errorf("command not found: %s", c.name))
	}
}

func (c commandNotFound) componentHandlers() map[string]ComponentHandler {
	return nil
}
