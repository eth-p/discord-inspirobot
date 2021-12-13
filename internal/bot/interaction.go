package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	interactionFlagEphemeral = 1 << 6 // Only the user can see the interaction.
)

// Interaction is a wrapper around the discordgo.Interaction struct.
// This provides some convenience features to make it easier to work with the API.
type Interaction struct {
	context     *botContext
	Interaction *discordgo.Interaction
}

// Respond sends a response to the interaction.
//
// Example:
//   i.Respond(&discordgo.InteractionResponseData {
//     Content: "Hey there! Nice to meet you!",
//   }
func (i Interaction) Respond(data *discordgo.InteractionResponseData) {
	if err := i.TryRespond(data); err != nil {
		exceptional{}.FatalGuard(fmt.Errorf("failed to respond to interaction: %w", err))
	}
}

// RespondPrivately sends a response to the interaction that only the user can see.
//
// Example:
//   i.RespondPrivately(&discordgo.InteractionResponseData {
//     Content: "Hey there! Nice to meet you!",
//   }
func (i Interaction) RespondPrivately(data *discordgo.InteractionResponseData) {
	if err := i.TryRespondPrivately(data); err != nil {
		exceptional{}.FatalGuard(fmt.Errorf("failed to respond to interaction: %w", err))
	}
}

// WillRespond informs the user that a response will be coming soon.
// This should be used if generating the response is expected to take longer than a second.
func (i Interaction) WillRespond() {
	if err := i.TryWillRespond(); err != nil {
		exceptional{}.FatalGuard(fmt.Errorf("failed to respond to interaction: %w", err))
	}
}

// TryRespond sends a response to the interaction, similarly to Respond.
// Rather than stopping execution, it may return an error.
func (i Interaction) TryRespond(data *discordgo.InteractionResponseData) error {
	return i.context.Bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// TryRespondPrivately sends a response that only the user can see, similarly to RespondPrivately.
// Rather than stopping execution, it may return an error.
func (i Interaction) TryRespondPrivately(data *discordgo.InteractionResponseData) error {
	newData := *data
	newData.Flags |= interactionFlagEphemeral
	return i.TryRespond(&newData)
}

// TryWillRespond informs the user that a response will be coming soon, similarly to Respond.
// Rather than stopping execution, it may return an error.
func (i Interaction) TryWillRespond() error {
	return i.context.Bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: interactionFlagEphemeral,
		},
	})
}
