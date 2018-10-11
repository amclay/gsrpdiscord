package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/goware/urlx"
	"github.com/nicklaw5/helix"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// To add the bot, have server owner go this URL
// https://discordapp.com/oauth2/authorize?client_id=498388551277215745&scope=bot

// Required Setup

// Add bot (above)
// Add bot to a "Bot" role in server
//  Allow this role access to change roles
//  Reorder this role above any roles that need changing (it can't alter roles of users that are above it)
// Create "streamer" role
//  Ensure this is below the bot role so it can alter users in this role

const (
	guildID        = "138425546173448192"
	streamerRoleID = "498416072563884042"
	botAPIToken    = "" // secret

	twitchClientID = "s4rhf3o6glc9dqh2vqizj77d7n2ztmx" // uses Xanbot, owned by user ID 39141793 (xangold)
)
const (
	updateInterval = 30 * time.Second
)

// need to be playing a game with this in title
var applicableStreamingTerms = []string{"ark", "gsrp", "gunsmoke", "chrome"}

var twitchClient *helix.Client

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	twitchClient, err = helix.NewClient(&helix.Options{
		ClientID: twitchClientID,
	})
	if err != nil {
		log.Fatalf("error creating twitch client: %v", err)
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", botAPIToken))
	if err != nil {
		log.Fatalf("error creating discord session: %v", err)
	}

	// TBD - not sure why this is required
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	discord.AddHandler(ready)
	discord.AddHandler(memberAdd)
	discord.AddHandler(guildMembersChunk)
	discord.AddHandler(memberRemove)
	discord.AddHandler(memberUpdate)

	loop(discord)
}

func loop(discord *discordgo.Session) {
	for {
		updateFromPresence(discord)
		time.Sleep(updateInterval)
	}
}

func updateFromPresence(discord *discordgo.Session) error {
	log.Println("starting main loop")
	var streamersMap = make(map[string]struct{})
	guild, err := discord.Guild(guildID)
	if err != nil {
		return fmt.Errorf("error getting guild: %v", err)
	}

	// get existing guild members, in order to find currently marked streamers
	var g errgroup.Group
	for _, gp := range guild.Presences {
		g.Go(func() error {
			guildMember, err := discord.GuildMember(guildID, gp.User.ID)
			if err != nil {
				return fmt.Errorf("error getting member: %v", err)
			}
			if stringInSlice(streamerRoleID, guildMember.Roles) {
				log.Printf("error adding streamer: %s", gp.User.ID)
				streamersMap[gp.User.ID] = struct{}{}
			}
			return nil
		})
	}
	g.Wait()

	for _, userPresence := range guild.Presences {
		isStreamingARK := false
		if userPresence.Game != nil {
			if len(userPresence.Game.URL) > 0 {
				hasRequiredTextInGame := false
				for _, term := range applicableStreamingTerms {
					if strings.Contains(strings.ToLower(userPresence.Game.Name), term) {
						hasRequiredTextInGame = true
					}
				}
				if hasRequiredTextInGame {
					log.Printf("%s playing %s on %s", userPresence.User.ID, userPresence.Game.Name, userPresence.Game.URL)
					if err := discord.GuildMemberRoleAdd(guildID, userPresence.User.ID, streamerRoleID); err != nil {
						return fmt.Errorf("could not add %s to role: %v", userPresence.User.ID, err)
					}
					isStreamingARK = true
				}
				title, err := getStreamTitle(userPresence.Game.URL)
				if err == nil {
					for _, term := range applicableStreamingTerms {
						if strings.Contains(strings.ToLower(title), term) {
							hasRequiredTextInGame = true
						}
					}
				}
			}
		}

		// remove the streamers that were previously playing ARK, but are no longer
		if !isStreamingARK {
			if _, ok := streamersMap[userPresence.User.ID]; ok {
				log.Printf("removing %s as a streamer...", userPresence.User.ID)
				if err := discord.GuildMemberRoleRemove(guildID, userPresence.User.ID, streamerRoleID); err != nil {
					return fmt.Errorf("could not remove %s from role: %v", userPresence.User.ID, err)
				}
			}
		}
	}
	log.Println("ending main loop")
	return nil
}

func getStreamTitle(url string) (string, error) {
	parsedURL, err := urlx.Parse(url)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %s %v", url, err)
	}
	twitchStreamResponse, err := twitchClient.GetStreams(&helix.StreamsParams{
		UserLogins: []string{parsedURL.Path},
	})
	if err != nil {
		return "", errors.Wrap(err, "could not get streams response")
	}
	if len(twitchStreamResponse.Data.Streams) == 0 {
		return "", errors.New("no streams")
	}
	return twitchStreamResponse.Data.Streams[0].Title, nil
}

// TBD - not sure what exactly these do? they seem to be required.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready", s, event)
}

func memberAdd(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready", s, event)
}

func guildMembersChunk(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready", s, event)
}

func memberRemove(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready", s, event)
}
func memberUpdate(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready", s, event)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
