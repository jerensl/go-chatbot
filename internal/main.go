package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := NewDiscordBotConnection()
	if err != nil {
		panic(err)
	}

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	discord.AddHandler(DiscordMessageHandle)
	discord.AddHandler(DiscordCommandHandle)
	discord.AddHandler(DiscordChannelHandle)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	registeredCommands := make(map[string]string, len(commands))
	for _, command := range commands {
		cmd, err := discord.ApplicationCommandCreate(os.Getenv("DISCORD_APP_ID"), os.Getenv("DISCORD_GUILD_ID"), command)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v \n", command.Name, err)
			return
		}
		registeredCommands[cmd.ID] = cmd.Name
	}

	fmt.Println("discord bot running...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	for id, name := range registeredCommands {
		err := discord.ApplicationCommandDelete(os.Getenv("DISCORD_APP_ID"), os.Getenv("DISCORD_GUILD_ID"), id)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v \n", name, err)
		}
	}
	discord.Close()
}