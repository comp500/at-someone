package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	modifiedContent := m.Content

	if strings.Contains(m.Content, "@someone") {
		modifiedContent = strings.Replace(m.Content, "@someone", "@someone", 1)
	}

	split := strings.Split(modifiedContent, " ")

	if split[0] == "@someone" {
		// Bot commands
		switch strings.ToLower(split[1]) {
		case "ping":
			// Respond to ping messages
			s.ChannelMessageSend(m.ChannelID, "help i've fallen and i can't get up i need @someone")
			return
		case "invite":
			// Send invite link
			s.ChannelMessageSend(m.ChannelID, "Invite link: https://discordapp.com/api/oauth2/authorize?client_id="+s.State.User.ID+"&permissions=2048&scope=bot")
			return
		}
	}

	if strings.Contains(modifiedContent, "@someone") {
		fmt.Print(modifiedContent)
		s.ChannelMessageSend(m.ChannelID, strings.Replace(modifiedContent, "@someone", "bot test", 1))
	}
}
