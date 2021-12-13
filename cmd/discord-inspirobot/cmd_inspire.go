package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	lru "github.com/hashicorp/golang-lru"

	"github.com/eth-p/discord-inspirobot/internal/bot"
)

type CmdInspire struct {
	bot.Command
}

func init() {
	usedResponses, _ := lru.New(1000)

	Commands = append(Commands, &CmdInspire{
		Command: bot.Command{
			Spec: &discordgo.ApplicationCommand{
				Name:        "inspire",
				Description: "Get an inspirational quote from InspiroBot.",
			},
			Handler: func(ctx *bot.CommandContext, interaction bot.Interaction) {
				err := respondWithQuote(interaction)
				if err != nil {
					ctx.Fatal(err)
				}
			},

			ComponentHandlers: map[string]bot.ComponentHandler{
				"another": func(ctx *bot.ComponentContext, interaction bot.Interaction) {
					err := respondWithQuote(interaction)
					if err != nil {
						ctx.Fatal(err)
					}
				},

				"share": func(ctx *bot.ComponentContext, interaction bot.Interaction) {
					// Parse the CustomData for the LRU cache key and the URL.
					split := strings.SplitN(ctx.CustomData, ",", 2)

					quoteUrl := split[1]
					titleId, err := strconv.Atoi(split[0])
					ctx.FatalGuard(err)

					// If the button was already pressed, don't share it again.
					if usedResponses.Contains(quoteUrl) {
						interaction.RespondPrivately(&discordgo.InteractionResponseData{
							Content: "You already shared this quote.",
						})
						return
					}

					// Respond with the contents of the original interaction.
					titleQuote := getInspirobotTitle(titleId)
					interaction.Respond(&discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{makeQuoteEmbed(titleQuote, quoteUrl)},
					})

					// Add it to the LRU to prevent spamming.
					usedResponses.Add(quoteUrl, nil)
				},
			},
		},
	})
}

// respondWithQuote responds to an interaction with an InspiroBot quote.
func respondWithQuote(interaction bot.Interaction) error {
	// If an inspirational quote could not be generated, tell the user and fail.
	url, err := generateInspirobotImage()
	if err != nil {
		interaction.RespondPrivately(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Unable to Inspire",
					Description: err.Error(),
					Color:       0xff0000,
				},
			},
		})

		return err
	}

	// Get the title quote.
	titleId, titleText := generateInspirobotTitle()

	// Share the inspirational quote.
	interaction.RespondPrivately(&discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{makeQuoteEmbed(titleText, url.String())},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Send to Channel",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("share:%v,%s", titleId, url.String()),
					},
					&discordgo.Button{
						Label:    "Another",
						Style:    discordgo.SecondaryButton,
						CustomID: "another",
					},
				},
			},
		},
	})

	return nil
}

// makeQuoteEmbed creates a discordgo.MessageEmbed instance for an InspiroBot quote.
func makeQuoteEmbed(title, url string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: title,
		Color: 0x2cc406,
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://inspirobot.me/",
			Name:    "InspiroBot",
			IconURL: "https://inspirobot.me/website/images/inspirobot-dark-green.png",
		},
	}
}
