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
				url, err := generateInspirobotImage()

				// If an inspirational quote could not be generated, tell the user and fail.
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

					ctx.Fatal(err)
				}

				// Get the title quote.
				titleId, titleQuote := generateInspirationalQuote()

				// Share the inspirational quote.
				interaction.RespondPrivately(&discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{makeQuoteEmbed(titleQuote, url.String())},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.Button{
									Label:    "Send to Channel",
									Style:    discordgo.PrimaryButton,
									CustomID: fmt.Sprintf("share:%v,%s", titleId, url.String()),
								},
							},
						},
					},
				})
			},

			ComponentHandlers: map[string]bot.ComponentHandler{
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
					titleQuote := getInspirationalQuote(titleId)
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

func makeQuoteEmbed(titleQuote, url string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: titleQuote,
		Color: 0x2cc406,
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
	}
}
