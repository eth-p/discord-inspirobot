package bot

// ComponentHandler is the function used to handle a component interaction.
type ComponentHandler func(ctx *ComponentContext, interaction Interaction)

// componentHandlerRegistration is a registration of a ComponentHandler that came from a Command.
// This is used to map Discord events to the appropriate handler.
type componentHandlerRegistration struct {
	Command CommandRegistration
	Handler ComponentHandler
}

// ---------------------------------------------------------------------------------------------------------------------

// ComponentContext is a context.Context passed to a ComponentHandler.
// This contains additional data about the bot and the command instance.
type ComponentContext struct {
	exceptional
	botContext

	ID         string
	CustomData string
	Values     []string

	// Command is the pointer to the "Command" struct which the ComponentHandler came from.
	// This may be used with type assertions to get the original struct back.
	//
	// Example:
	//   func MyHandler(ctx ComponentContext, interaction bot.Interaction) error {
	//     self := ctx.Command.(*MyCommand)
	//     self.AnyCustomField = 3
	//   }
	Command interface{}
}
