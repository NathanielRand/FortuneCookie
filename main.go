package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var version string = "1.0.0"

func goDotEnvVariable(key string) string {
	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return value from key provided.
	return os.Getenv(key)
}

func main() {
	// Grab bot token env var.
	botToken := goDotEnvVariable("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	if strings.Contains(content, "!vthelp") {
		// Build help message
		author := m.Author.Username

		// Title
		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"

		// Notes
		note1 := "- Bot will return a fortune based on unfathomable cosmic events. \n"
		note2 := "- Commands are case-sensitive. They must be in lower-case :) \n"
		note3 := "- Dev: Narsiq#5638. DM me for requests/questions/love. \n"

		// Commands
		commandHelp := "❔  !fchelp : Provides a list of my commands. \n"
		commandFortune := "🦶🏽  !fc : Return a fortune based on unfathomable cosmic events. \n"
		commandInvite := "🔗  !fcinvite : A invite link for the FortuneCookie Bot. \n"
		commandSite := "🔗  !fcsite : Link to the FortuneCookie website. \n"
		commandSupport := "✨  !fcsupport : Link to the FortuneCookie Patreon. \n"
		commandStats := "📊  !fcstats : Check out FortuneCookie stats. \n"
		commandVersion := "🤖  !fcversion : Current FortuneCookie version. \n"

		// Build sub messages
		notesMessage := note1 + note2 + note3
		commandsMessage := commandHelp + commandFortune
		othersMessage := commandInvite + commandSite + commandSupport + commandStats + commandVersion

		// Build full message
		message := "Whats up " + author + "\n \n" + commandHelpTitle + "NOTES: \n \n" + notesMessage + "\n" + "COMMANDS: \n \n" + commandsMessage + "\n" + "OTHER: \n \n" + othersMessage + "\n \n" + "https://www.patreon.com/BotVoteTo"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!fcsite") {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!fcsupport") {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!fcversion") {
		// Build start vote message
		message := "FortuneCookie is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!fcstats") {
		// TODO: This will need to be updated to iterate through
		// all shards once the bot joins 1,000 servers.
		guilds := s.State.Ready.Guilds
		fmt.Println(len(guilds))
		guildCount := len(guilds)

		guildCountStr := strconv.Itoa(guildCount)

		// // Build start vote message
		message := "FortuneCookie is currently on " + guildCountStr + " servers. Such wow!"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!fcinvite") {
		author := m.Author.Username

		// // Build start vote message
		message := "Wow! Such nice " + author + ". Thanks for spreading the 💖. Here is an invite link made just for you... \n \n" + "https://discord.com/api/oauth2/authorize?client_id=921252848036106270&permissions=274877995072&scope=bot"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.EqualFold(content, "!fc") {
		fortune := getFortune()

		// Grab author
		author := m.Author.Username

		// Build start vote message
		message := "🥠 Yo " + author + "... " + fortune

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type Fortune struct {
	Message string
}

func getFortune() string {
	csvFile, err := os.Open("fortunes.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	// Read csv file
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	// Generate random number using min/max index of csv file lines.
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 858
	randomIndex := rand.Intn(max-min+1) + min
	result := csvLines[randomIndex]

	return result[0]
}
