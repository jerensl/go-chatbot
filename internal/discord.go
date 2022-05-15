package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	*discordgo.Session
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = true
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "google-it",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "say-hello",
			Type: discordgo.UserApplicationCommand,
		},
	}
)

func NewDiscordBotConnection() (*DiscordBot, error) {
	discord, err := discordgo.New("Bot "+os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, errors.New("failed to create new connection discord")
	}
	
	return &DiscordBot{
		discord,
	}, nil
}


func DiscordChannelHandle(s *discordgo.Session, channel *discordgo.ChannelCreate) {
	if channel.Type == discordgo.ChannelTypeDM {
		fmt.Println(channel.Messages)
	}
}

func DiscordMessageHandle(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	} 

	if msg.Type == discordgo.MessageType(discordgo.ChannelTypeDM) {
		fmt.Printf("DM Mesage %v \n", msg.Content)
	}


	privateChannel, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		fmt.Println("Cannot create new channel err", err)
		return
	}

	user, _ := s.User(msg.Author.ID) 


	_, err = s.ChannelMessageSend(privateChannel.ID, "Hello "+user.Username)
	if err != nil {
		fmt.Println("Failed to send message channel")
		return
	}
}

func searchLink(message, format, sep string) string {
	return fmt.Sprintf(format, strings.Join(
		strings.Split(
			message,
			" ",
		),
		sep,
	))
}


func DiscordCommandHandle(s *discordgo.Session, i *discordgo.InteractionCreate)  {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	switch data.Name {
	case "say-hallo":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Hi %v, nice to meet you", i.User.Username),
			},
		})
	case "google-it":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: searchLink(
					i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content,
					"https://google.com/search?q=%s", "+"),
				Flags: 1 << 6,
			},
		})
	}
}