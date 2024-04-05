package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Uwu is a handler that sends a message when a message contains uwu
func Uwu(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "uwu") {
		span := tracer.StartSpan(
			"commands.handlers.uwu:Uwu",
			tracer.ResourceName("Handlers.Uwu"),
		)
		defer span.Finish()
		message, err := s.ChannelMessageSend(m.ChannelID, "**𝓤𝔀𝓤**")
		if err != nil {
			logging.Error(s, err.Error(), m.Member.User, span, logrus.Fields{"error": err})
		} else {
			logging.DebugButton(
				s,
				"**𝓤𝔀𝓤**\n",
				discordgo.Button{
					Label: "View Message",
					URL:   helpers.JumpURL(message),
					Style: discordgo.LinkButton,
					Emoji: &discordgo.ComponentEmoji{
						Name: "👀",
					},
				},
				m.Member.User,
				span,
			)
		}
	}
}
