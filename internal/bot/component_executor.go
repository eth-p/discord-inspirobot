package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// executeComponent executes a ComponentHandler and returns its CommandContext and ExecutionStatus.
func executeComponent(ctx botContext, registration componentHandlerRegistration, customID string, customData string, data *discordgo.MessageComponentInteractionData, interaction *discordgo.Interaction) (*ComponentContext, ExecutionStatus) {
	compCtx := &ComponentContext{
		botContext: ctx,
		Command:    registration.Command,

		ID:         customID,
		CustomData: customData,
		Values:     data.Values,
	}

	// Run the component handler.
	status := ExecutionStatus{}
	status.TimeStarted = time.Now()
	status.Err = compCtx.try(func() {
		registration.Handler(compCtx, Interaction{
			context:     &compCtx.botContext,
			Interaction: interaction,
		})
	})
	status.TimeFinished = time.Now()

	// Return the context and status.
	return compCtx, status
}
