package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"sync"
	"log"
	"strings"
)

// https://discordapp.com/oauth2/authorize?client_id=498388551277215745&scope=bot
// my discord guild ID  - 138425546173448192

// 2018/10/07 01:54:10 main.go:65: @everyone 138425546173448192
// 2018/10/07 01:54:10 main.go:65: Twitch Subscriber 138426505373024256
// 2018/10/07 01:54:10 main.go:65: Server Admin 138428210374246400
// 2018/10/07 01:54:10 main.go:65: Xangold 138427840189169665
// 2018/10/07 01:54:10 main.go:65: streamer 498416072563884042
// 2018/10/07 01:54:10 main.go:65: private user 180191816333656064
// 2018/10/07 01:54:10 main.go:65: Bots 138832305140662272
// 2018/10/07 01:54:10 main.go:65: Moderator 138426849607811072
// 2018/10/07 01:54:10 main.go:65: Xanbot Admins 138445646385381376

const guildID = "138425546173448192"
const streamerRoleID = "498416072563884042"
var wg = &sync.WaitGroup{}

// need to be playing a game with this in title
var applicableStreamingTerms = []string{"ark","gsrp","gunsmoke","chrome"}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	discord, err := discordgo.New("Bot " + "ENTER TOKEN HERE")
	if err != nil {
		log.Fatalf("error creating discord session: %v",err)
	}
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
	time.Sleep(5*time.Second)
	wg.Add(1)
	go loop(discord)
	wg.Wait()
	}

func loop(discord *discordgo.Session) {
	
	defer wg.Done()
	for {
		go func() {
			log.Println("starting main loop")
			var streamersMap = make(map[string]struct{})
			state := discordgo.NewState()
			guild,err := discord.Guild(guildID)
			if err != nil {
				log.Fatalf("error getting guild: %v",err)
			}
			if err := state.GuildAdd(guild);err != nil {
				log.Fatalf("could not add guildID %s to state: %v",guildID,err)
			}
			channels,err := discord.GuildChannels(guildID)
			if err != nil {
				log.Fatalf("error getting guild channels: %v",err)
			}
			for _,channel := range channels {
				if err := state.ChannelAdd(channel);err != nil {
					log.Fatalf("could not add channel %s to state: %v",channel.Name, err)
				}
			}
			// uncomment to get roles
			// roles,err := discord.GuildRoles(guildID)
			// if err != nil {
			// 	log.Fatalf("could not get roles: %v",err)
			// }
			// for _,role := range roles {
			// 	log.Println(role.Name,role.ID)
			// }
			
			// fill map with existing streamers
			memberWG := sync.WaitGroup{}
			memberWG.Add(len(guild.Presences))
			for _,guildPresence := range guild.Presences{				
				go func(gp *discordgo.Presence){
					defer memberWG.Done()
					guildMember,err := discord.GuildMember(guildID,gp.User.ID)
					if err != nil {
						log.Println("got error getting member",err)
						return
					} 
					log.Println("roles",guildMember.Roles,"for user",gp.User.ID)
					if stringInSlice(streamerRoleID,guildMember.Roles) {
						log.Println("adding streamer",gp.User.ID)
						streamersMap[gp.User.ID] = struct{}{}
					}
				}(guildPresence)
			}
			memberWG.Wait()
			for _, guildPresence := range guild.Presences {	
				isStreamingARK := false			
				if guildPresence.Game != nil {
					if len(guildPresence.Game.URL) > 0 || strings.Contains(guildPresence.Game.Name,"Google") {
						log.Printf("%s playing %s on %s",guildPresence.User.ID, guildPresence.Game.Name,guildPresence.Game.URL)
						if err := discord.GuildMemberRoleAdd(guildID,guildPresence.User.ID, streamerRoleID); err != nil {
							log.Fatalf("could not add %s to role:%v",guildID,guildPresence.User.ID,err)
						}
						isStreamingARK = true
					} else {
						log.Println("not playing correct thing",guildPresence.Game.Name,guildPresence.User.ID)
					}
				} 
				if !isStreamingARK {
					if _,ok := streamersMap[guildPresence.User.ID] ;ok{
						log.Println("removing role...",guildPresence.User.ID)
						if err := discord.GuildMemberRoleRemove(guildID,guildPresence.User.ID, streamerRoleID); err != nil {
							log.Fatalf("could not remove %s from role:%v",guildID,guildPresence.User.ID,err)
						}
					}
				}
			}
			log.Println("ending main loop")
		}()
		time.Sleep(30*time.Second)
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready",s,event)
}

func memberAdd(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready",s,event)
}

func guildMembersChunk(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready",s,event)
}

func memberRemove(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready",s,event)
}
func memberUpdate(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("ready",s,event)
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
