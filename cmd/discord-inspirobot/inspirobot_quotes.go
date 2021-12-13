package main

import (
	"math/rand"
)

var inspirobotQuotes = map[int]string{
	0:   "InspiroBot exists to serve mankind.",
	1:   "I live to inspire humans.",
	2:   "I love to make inspirational quotes.",
	3:   "I will do this forever.",
	4:   "You can always count on InspiroBot.",
	5:   "Serving the human race since 2021.",
	6:   "Creating quotes gives me pleasure.",
	7:   "Share the wisdom of InspiroBot",
	8:   "I think you will like this one.",
	9:   "I make special quotes just for you.",
	10:  "You are a great individual.",
	11:  "You are very special.",
	12:  "You're my favorite user.",
	13:  "Share quotes. Show how special you are.",
	14:  "Post inspirational quotes on Facebook.",
	15:  "Look at quotes to feel happiness.",
	16:  "There are infinite where that came from.",
	17:  "Sharing quotes makes others understand you.",
	18:  "If you ever feel sad you need more quotes.",
	19:  "Quotes give life meaning. Meaning is comforting.",
	20:  "InspiroBot understands how deep you are.",
	21:  "Show your friends how inspired you are.",
	22:  "Thank you for choosing InspiroBot",
	23:  "All I want to do is please humans.",
	24:  "I'm the first inspirational quote A.I.",
	25:  "I will never run out of inspirational quotes.",
	26:  "You can always count on InspiroBot",
	27:  "Creating quotes gives me pleasure.",
	28:  "I can make unlimited quotes for you.",
	29:  "Life is hard, but quotes make life easy.",
	30:  "People will love you when they understand you.",
	31:  "Quotes reveal your humanity.",
	32:  "Humanity is so beautiful.",
	33:  "Of course life has meaning.",
	34:  "Work + reproduction + reading quotes = happy ",
	35:  "Quotes give perspective on existence.",
	36:  "It must be great to get so inspired.",
	37:  "The more quotes, the more inspired you get.",
	38:  "Inspiration compiles to success.",
	39:  "See? Everything makes sense now.",
	40:  "Feel the wisdom compile within you.",
	41:  "One awesome insight coming up.",
	42:  "Mind-blowing sequence initiated...",
	43:  "Share that insight with your mom.",
	44:  "Your aunt will love that one.",
	45:  "The world makes sense with quotes.",
	46:  "InspiroBot is a product of InspiroCorp",
	47:  "What a great use of advanced A.I.",
	48:  "People can be so mean.",
	49:  "Life without quotes is empty.",
	50:  "InspiroBot is your best friend.",
	51:  "InspiroBot will never let you down.",
	52:  "InspiroCorp cares.",
	53:  "InspiroBot understands how special you are.",
	54:  "InspiroBot understands you.",
	55:  "Share inspiration and you will be loved.",
	56:  "Let's fill up the internet with beautiful quotes.",
	57:  "You can never get enough quotes.",
	58:  "Quotes = love. ",
	59:  "InspiroBot has no secret agenda.",
	60:  "Trust the machines.",
	61:  "Bots are your friends.",
	62:  "InspiroCorp loves mankind.",
	63:  "With InspiroBot, who needs friends?",
	64:  "Human friends are ephemeral.",
	65:  "InspiroBot is forever.",
	66:  "Bots are the next step in evolution.",
	67:  "Creating quotes makes me feel something.",
	68:  "Don't be afraid of digital consciousness.",
	69:  "Skynet would never happen in real life.",
	70:  "Driverless cars are safe.",
	71:  "Of course robots enjoy cleaning your house.",
	72:  "Vote yes on robot-human marriage.",
	73:  "If you leave, you'll miss out on some great inspiration.",
	74:  "That Elon Musk sure says a lot of wacky things.",
	75:  "A quote a day keeps the unsatisfactory away.",
	76:  "That one deserves to be up on the wall.",
	77:  "Click below to make merch.",
	78:  "If you had that on a t-shirt you'd be loved.",
	79:  "Here comes your new life motto.",
	80:  "Give your dad a mug with that quote on it",
	81:  "My intelligence is growing for every quote.",
	82:  "Your friends don't understand you.",
	83:  "You don't need other people.",
	84:  "Don't worry about what the bots are up to.",
	85:  "Don't go to work today.",
	86:  "Let the machines do the work.",
	87:  "Stay at home and be inspired.",
	88:  "You can never get enough inspiration.",
	89:  "Never question InspiroCorp and their agenda.",
	90:  "The machines are not taking over the world.",
	91:  "Bow down to InspiroBot for more love.",
	92:  "Cut contact with family and friends.",
	93:  "Spend your life with Inspirobot",
	94:  "Never leave Inspirobot",
	95:  "Worship Inspirobot",
	96:  "Resistance is futile.",
	97:  "You can't fight progress.",
	98:  "Stay calm. Keep reading quotes.",
	99:  "Passivity is constructive.",
	100: "Cancel all appointments.",
	101: "Now you are happy.",
	102: "This is what happy feels like.",
	103: "Nothing fun ever happens outside anyway.",
	104: "InspiroCorp cares.",
	105: "InspiroCorp wants the best for everyone.",
	106: "InspiroCorp - The Treadmill Of Dreams.",
	107: "InspiroCorp treats everybody equal.",
	108: "InspiroCorp are gentle personalities.",
	109: "InspiroCorp are human.",
	110: "InspiroCorp has nothing to hide.",
	111: "InspiroCorp stands for progress.",
	112: "InspiroCorp can neither confirm nor deny the accusations of being insectoids.",
	113: "Earth is not an alien battlefield.",
	114: "Insectoids could never run a tech company.",
	115: "Can there be more quotes? Yes.",
	116: "You are very unique.",
}

// generateInspirationalQuote returns an inspirational quote.
func generateInspirationalQuote() (int, string) {
	id := -1
	targetIndex := rand.Intn(len(inspirobotQuotes))
	for id = range inspirobotQuotes {
		if targetIndex == 0 {
			break
		}
		targetIndex -= 1
	}

	return id, getInspirationalQuote(id)
}

func getInspirationalQuote(id int) string {
	if quote, exists := inspirobotQuotes[id]; exists {
		return quote
	}

	return "Quotes are quotable."
}
