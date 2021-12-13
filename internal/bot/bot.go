package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog/v2"
)

// Bot is the main instance of the discord-inspirobot.
type Bot struct {
	Session *discordgo.Session
	appID   string

	commands          map[string]CommandRegistration
	componentHandlers map[string]componentHandlerRegistration

	ctx       context.Context
	stop      chan struct{}
	stopError error
}

// uninitializedBot is an unexported type that represents the Bot before it has connected to the Discord Gateway.
// This is provided to Option instances to set up the bot.
type uninitializedBot *Bot

// botContext is a context that contains a pointer to the active Bot.
// This is intended to be composed inside of other context structs.
type botContext struct {
	context.Context
	Bot *Bot
}

// ---------------------------------------------------------------------------------------------------------------------

// Wait waits for the bot to stop.
func (b *Bot) Wait() error {
	<-b.stop
	return b.stopError
}

// AppID returns the bot's Application ID.
func (b *Bot) AppID() string {
	return b.appID
}

// ---------------------------------------------------------------------------------------------------------------------

// New connects to Discord and starts the Bot instance.
func New(token string, ctx context.Context, options ...Option) (*Bot, error) {
	// Create the session.
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Create the bot instance.
	bot := &Bot{
		Session: session,

		commands: make(map[string]CommandRegistration),

		ctx:       ctx,
		stopError: nil,
		stop:      make(chan struct{}),
	}

	// Initialize the bot instance.
	klog.V(4).InfoS("Initializing bot...")
	for _, option := range options {
		if err := option(uninitializedBot(bot)); err != nil {
			return nil, fmt.Errorf("unable to set up option: %w", err)
		}
	}

	// Create a channel to wait upon for ready.
	ready := make(chan struct{})
	session.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		ready <- struct{}{}
	})

	// Start the bot.
	klog.V(4).InfoS("Opening Discord session...")
	if err = session.Open(); err != nil {
		return nil, err
	}

	// Start a goroutine to stop the bot when the context expires.
	go func() {
		<-ctx.Done()
		if err := bot.shutdown(); err != nil {
			bot.stopError = err
			klog.V(1).ErrorS(err, "Unable to close Discord session.")
		}

		close(bot.stop)
	}()

	// Wait for the bot to become ready.
	klog.V(4).InfoS("Waiting for Discord Ready...")
	select {
	case <-ready: // Done
	case <-ctx.Done(): // Cancelled
	}

	close(ready)

	// Return the context error (if expired).
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Set up the bot.
	bot.appID = bot.Session.State.User.ID
	for _, setupFunc := range []func() error{
		bot.setupCommandHandlers,
		bot.setupComponentHandlers,
		bot.setupCommandRegistrations,
	} {
		if err = setupFunc(); err != nil {
			_ = bot.shutdown()
			return nil, err
		}
	}

	// Return the bot instance.
	klog.V(4).InfoS("Bot ready.")
	return bot, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// shutdown shuts down the bot.
func (b *Bot) shutdown() error {
	klog.V(1).InfoS("Shutting down bot.")

	// Disconnect from the gateway.
	if err := b.Session.Close(); err != nil {
		return err
	}

	return nil
}

// newExecutionContext creates a new botContext for an execution.
func (b *Bot) newExecutionContext() botContext {
	return botContext{Bot: b, Context: b.ctx}
}

// setupCommandHandlers sets up a Discord event listener which handles slash commands and delegates control flow to the
// appropriate CommandHandler instance.
func (b *Bot) setupCommandHandlers() error {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		// Get the command handler.
		// If the command wasn't found, use a fake command that returns an error.
		command, ok := b.commands[i.ApplicationCommandData().Name]
		if !ok {
			command = commandNotFound{name: i.ApplicationCommandData().Name}
		}

		// Run the command.
		go func() {
			_, status := executeCommand(b.newExecutionContext(), command, i.Interaction)
			statusV := 3
			if status.Err != nil {
				statusV = 2
			}

			// Log the command result.
			klog.V(klog.Level(statusV)).InfoS("Received command from user.",
				"command", i.ApplicationCommandData().Name,
				"user", i.Member.User.ID,
				"duration", status.Duration(),
				"err", status.Err,
			)
		}()
	})

	return nil
}

// setupCommandHandlers sets up a Discord event listener which handles interaction components and delegates control flow
// to the appropriate CommandHandler instance.
func (b *Bot) setupComponentHandlers() error {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		// Get the component data.
		componentData := i.MessageComponentData()
		customID := componentData.CustomID
		customData := ""

		if strings.Contains(customID, ":") {
			split := strings.SplitN(customID, ":", 2)
			customID = split[0]
			customData = split[1]
		}

		// Get the component handler.
		// If the handler wasn't found, use a fake handler that returns an error.
		registration, ok := b.componentHandlers[customID]
		if !ok {
			registration = componentHandlerRegistration{
				Command: commandNotFound{name: "?"},
				Handler: func(ctx *ComponentContext, interaction Interaction) {
					ctx.Fatal(fmt.Errorf("component handler not found: %s", ctx.ID))
				},
			}
		}

		// Run the handler.
		go func() {
			_, status := executeComponent(b.newExecutionContext(), registration, customID, customData, &componentData, i.Interaction)
			statusV := 3
			if status.Err != nil {
				statusV = 2
			}

			// Log the command result.
			klog.V(klog.Level(statusV)).InfoS("Received component interaction from user.",
				"component", customID,
				"user", i.Member.User.ID,
				"duration", status.Duration(),
				"err", status.Err,
			)
		}()
	})

	return nil
}

// setupCommandRegistrations informs the Discord Gateway of the slash commands that this bot provides.
// This registers the commands for all guilds.
func (b *Bot) setupCommandRegistrations() error {
	commandSpecs := make([]*discordgo.ApplicationCommand, 0, len(b.commands))
	for _, command := range b.commands {
		commandSpecs = append(commandSpecs, command.spec())
	}

	// Overwrite all the existing commands.
	_, err := b.Session.ApplicationCommandBulkOverwrite(b.appID, "", commandSpecs)
	if err != nil {
		return fmt.Errorf("unable to register slash commands: %w", err)
	}

	// Regenerate the component handler registrations.
	b.componentHandlers = make(map[string]componentHandlerRegistration)
	for _, command := range b.commands {
		for handlerId, handlerFunc := range command.componentHandlers() {
			if _, exists := b.componentHandlers[handlerId]; exists {
				return fmt.Errorf("component handler already exists for ID '%s'", handlerId)
			}

			b.componentHandlers[handlerId] = componentHandlerRegistration{
				Command: command,
				Handler: handlerFunc,
			}
		}
	}

	klog.V(3).InfoS("Registered slash commands.", "count", len(b.commands))
	return nil
}
