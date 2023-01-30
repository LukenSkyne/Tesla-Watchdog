package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"math"
	"time"
)

const (
	Green  = 0x3ba55c
	Yellow = 0xd8b576
	Red    = 0xdb5368
)

type Discord struct {
	log       *zap.SugaredLogger
	session   *discordgo.Session
	botToken  string
	channelId string
}

func NewDiscord(log *zap.SugaredLogger, botToken string, channelId string) *Discord {
	return &Discord{
		log:       log,
		botToken:  botToken,
		channelId: channelId,
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

	go d.sendStartupMessage()

	return d.log.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zap.WarnLevel || entry.Level == zap.ErrorLevel {
			go d.onLogCallback(entry)
		}

		return nil
	})), true
}

func (d *Discord) sendStartupMessage() {
	tmp := &discordgo.MessageSend{
		Content: "",
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

func (d *Discord) onLogCallback(entry zapcore.Entry) {
	var color int

	switch entry.Level {
	case zap.WarnLevel:
		color = Yellow // yellow
	case zap.ErrorLevel:
		color = Red // red
	default:
		color = 0
	}

	length := int(math.Min(256, float64(len(entry.Message))))
	tmp := &discordgo.MessageSend{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       entry.Level.CapitalString() + " @ " + entry.Caller.TrimmedPath(),
				Description: entry.Message[:length],
				Timestamp:   time.Now().Format(time.RFC3339),
				Color:       color,
			},
		},
	}

	if _, err := d.session.ChannelMessageSendComplex(d.channelId, tmp); err != nil {
		d.log.Debug("discord error: %v\n", err)
	}
}
