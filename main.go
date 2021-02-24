package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"wacha/utils"
)

func init() {
	if err := gotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// regexps
	r, err := regexp.Compile("!wacha\\s([\\w]+)\\s([\\w]+)")
	if err != nil {
		log.Fatalln(err)
	}

	if m.Author.ID == s.State.User.ID {
		return
	}
	if r.MatchString(m.Content) {
		d := utils.GetParams("!wacha\\s(?P<mode>[\\w]+)\\s(?P<option>[\\w]+)", m.Content)
		switch d["mode"] {
		case "ban":
			if len(m.Mentions) > 0 && utils.CanBanMembers(m.Author) {
				banMember := m.Mentions[0]
				reason := strings.Replace(d["option"], "_", " ", -1)
				if os.Getenv("FakeBan") == "false" || os.Getenv("FakeBan") == "" {
					err = s.GuildBanCreate(m.GuildID, banMember.ID, 7)
				}
				if err != nil {
					_, _ = s.ChannelMessageSend(m.ChannelID, "cant ban "+banMember.Username)
					fmt.Println(err.Error())
					return
				}
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Banned %s Reason: %s (This will be logged on the log channel)", banMember.Username, reason))
				_, err = s.ChannelMessageSend(os.Getenv("BanLogChannelId"), fmt.Sprintf("Banned %s Reason: %s on %v", banMember.Username, reason, time.Now()))
				if err != nil {
					log.Println("warning: " + err.Error())
				}
			} else {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "please mention a user to ban or you might not have the required role", m.MessageReference)
			}
		case "report":
			body := strings.Replace(d["option"], "_", " ", -1)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, os.Getenv("ReportSentMessage"), m.MessageReference)
			_, _ = s.ChannelMessageSend(os.Getenv("ReportLogChannel"), fmt.Sprintf("report by %s#%s on %v: %s", m.Author.Username, m.Author.Discriminator, time.Now(), body))
		}
	}
}

func main() {

	println(os.Getenv("ClientId"))

	dg, err := discordgo.New("Bot " + os.Getenv("AuthCode"))
	if err != nil {
		log.Fatalln(err)
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

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
