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

const prefix string = "!fc"
const version string = "1.1.0"

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

	if strings.Contains(content, prefix+"help") {
		// Build help message
		author := m.Author.Username

		// Title
		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"

		// Notes
		note1 := "- Bot will return a fortune based on unfathomable cosmic events. \n"
		note2 := "- Commands are case-sensitive. They must be in lower-case :) \n"
		note3 := "- Dev: Narsiq#5638. DM me for requests/questions/love. \n"

		// Commands
		commandHelp := "❔  " + prefix + "help : Provides a list of my commands. \n"
		commandFortune := "🦶🏽  " + prefix + "!fc : Return a fortune based on unfathomable cosmic events. \n"
		commandInvite := "🔗  " + prefix + "invite : A invite link for the FortuneCookie Bot. \n"
		commandSite := "🔗  " + prefix + "site : Link to the FortuneCookie website. \n"
		commandSupport := "✨  " + prefix + "support : Link to the FortuneCookie Patreon. \n"
		commandStats := "📊  " + prefix + "stats : Check out FortuneCookie stats. \n"
		commandVersion := "🤖  " + prefix + "version : Current FortuneCookie version. \n"

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

	if strings.Contains(content, prefix+"site") {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"support") {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"version") {
		// Build start vote message
		message := "FortuneCookie is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"stats") {
		// TODO: This will need to be updated to iterate through
		// all shards once the bot joins 1,000 servers.
		guilds := s.State.Ready.Guilds
		fmt.Println(len(guilds))
		guildCount := len(guilds)

		guildCountStr := strconv.Itoa(guildCount)

		// // Build start vote message
		message := "FortuneCookie is currently on " + guildCountStr + " servers. Such wow!"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.EqualFold(content, prefix+"invite") {
		author := m.Author.Username

		// // Build start vote message
		message := "Wow! Such nice " + author + ". Thanks for spreading the 💖. Here is an invite link made just for you... \n \n" + "https://discord.com/api/oauth2/authorize?client_id=921252848036106270&permissions=274877995072&scope=bot"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.EqualFold(content, prefix) {
		// Call greeting, fortune, and five numbers funcs to generate values.
		greeting := getGreeting()
		fortune := getFortune()
		fiveNumbers := getFiveNumbers()

		// Convert map values into string to be used in message below.
		one := strconv.Itoa(fiveNumbers[1])
		two := strconv.Itoa(fiveNumbers[2])
		three := strconv.Itoa(fiveNumbers[3])
		four := strconv.Itoa(fiveNumbers[4])
		five := strconv.Itoa(fiveNumbers[5])
		six := strconv.Itoa(getOneNumber())

		// Grab author
		author := m.Author.Username

		// Build start vote message
		messageGreet := greeting + " " + author + "... \n"
		messageFortune := "```fix" + "\n" + "🥠 " + fortune + "\n" + "```"
		messageTitle := "🍀 " + author + "'s Lucky Numbers: "
		messageLucky := "```CSS" + "\n" + messageTitle + one + "-" + two + "-" + three + "-" + four + "-" + five + "-[" + six + "]" + "\n" + "```"
		messageFull := messageGreet + messageFortune + messageLucky

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, messageFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}
}

type Greeting struct {
	Message string
}

func getGreeting() string {
	csvFile, err := os.Open("greetings.csv")
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
	max := 105
	randomIndex := rand.Intn(max-min+1) + min
	result := csvLines[randomIndex]

	return result[0]
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
	max := 856
	randomIndex := rand.Intn(max-min+1) + min
	result := csvLines[randomIndex]

	return result[0]
}

func getFiveNumbers() map[int]int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 70

	m := make(map[int]int)
	m[1] = rand.Intn(max-min) + min
	m[2] = rand.Intn(max-min) + min
	m[3] = rand.Intn(max-min) + min
	m[4] = rand.Intn(max-min) + min
	m[5] = rand.Intn(max-min) + min

	for hasDupes(m) {
		m[1] = rand.Intn(max-min) + min
		m[2] = rand.Intn(max-min) + min
		m[3] = rand.Intn(max-min) + min
		m[4] = rand.Intn(max-min) + min
		m[5] = rand.Intn(max-min) + min
	}

	return m
}

func getOneNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 25
	result := rand.Intn(max-min) + min

	return result
}

func hasDupes(m map[int]int) bool {
	x := make(map[int]struct{})
	for _, v := range m {
		if _, has := x[v]; has {
			return true
		}
		x[v] = struct{}{}
	}
	return false
}
