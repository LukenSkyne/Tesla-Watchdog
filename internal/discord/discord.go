package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

const (
	Green  = 0x3ba55c
	Yellow = 0xd8b576
	Red    = 0xdb5368
)

type Discord struct {
	session   *discordgo.Session
	log       *zap.SugaredLogger
	botToken  string
	channelId string
	userId    string
}

func NewDiscord(log *zap.SugaredLogger) *Discord {
	return &Discord{
		session:   nil,
		log:       log,
		botToken:  os.Getenv("BOT_TOKEN"),
		channelId: os.Getenv("CHANNEL_ID"),
		userId:    os.Getenv("USER_ID"),
	}
}

func (d *Discord) Stop() {
	_ = d.session.Close()
}

func (d *Discord) Start() (*zap.SugaredLogger, bool) {
	if d.botToken == "" {
		d.log.Warn("no Token for Discord Bot supplied")
		return d.log, false
	}

	var err error
	d.session, err = discordgo.New("Bot " + d.botToken)

	if err != nil {
		log.Fatalf("failed to create Bot session: %v\n", err)
	}

	err = d.session.Open()

	if err != nil {
		log.Fatalf("cannot open Bot session: %v\n", err)
	}

	d.sendStartupMessage()

	return d.log.WithOptions(zap.Hooks(d.onLogCallback)), true
}

func (d *Discord) sendStartupMessage() {
	tmp := &discordgo.MessageSend{
		Content: "<@" + d.userId + ">",
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:      discordgo.EmbedTypeRich,
				Title:     "Watchdog Started",
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     Green,
			},
		},
	}

	if _, err := d.session.ChannelMessageSendComplex(d.channelId, tmp); err != nil {
		d.log.Debug("discord error: %v\n", err)
	}
}

func (d *Discord) onLogCallback(entry zapcore.Entry) error {
	if entry.Level == zap.InfoLevel || entry.Level == zap.DebugLevel {
		return nil
	}

	var color int

	switch entry.Level {
	case zap.WarnLevel:
		color = Yellow // yellow
	case zap.ErrorLevel:
		color = Red // red
	default:
		color = 0
	}

	tmp := &discordgo.MessageSend{
		Content: "<@" + d.userId + ">",
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       entry.Level.CapitalString() + " @ " + entry.Caller.TrimmedPath(),
				Description: entry.Message,
				Timestamp:   time.Now().Format(time.RFC3339),
				Color:       color,
			},
		},
	}

	if _, err := d.session.ChannelMessageSendComplex(d.channelId, tmp); err != nil {
		d.log.Debug("discord error: %v\n", err)
	}

	return nil
}
