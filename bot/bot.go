package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	env "github.com/DRBC17/solaire-bot-go/plugins"
	"github.com/bwmarrin/discordgo"
)

var Token string = env.Get("TOKEN")
var s *discordgo.Session

func init() {
	var err error
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "registrar",
			Description: "Registrar miembro",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "discord_name",
					Description: "Nombre en el servidor @name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "albion_name",
					Description: "Nombre del juego",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"registrar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			msgformat := ""
			msgformat += "Sol Exánime,\n%s \n\n Comando: %s"

			albion_name := optionMap["albion_name"]
			discord_member_id := optionMap["discord_name"]
			commandHistory := ""

			if discord_member_id != nil && albion_name != nil {
				commandHistory = "/registrar discord_name: <@" + discord_member_id.Value.(string) + "> albion_name: " + albion_name.StringValue()
			}

			if discord_member_id != nil {

				api_url := env.Get("ALBION_API") + env.Get("ALBION_GUILD_ID") + "/members"

				resp, err := http.Get(api_url)

				if err != nil {
					log.Fatalf("Cannot get members: %v", err)
				}

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					log.Fatalf("Cannot read body: %v", err)
				}

				var members []Member
				var memberFound Member
				var memberFoundExist bool = false

				if err := json.Unmarshal(body, &members); err != nil {
					log.Fatalf("Cannot unmarshal members: %v", err)
				}

				for _, member := range members {
					if member.Name == albion_name.StringValue() {
						memberFound = member
						memberFoundExist = true
						break
					}
				}

				if memberFoundExist {
					if existMemberRole(s, env.Get("GUILD_ID"), discord_member_id.Value.(string), env.Get("MEMBER_ROLE_ID")) {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							// Ignore type for now, they will be discussed in "responses"
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf(
									msgformat,
									albion_name.StringValue()+" ya es miembro de la guild.",
									commandHistory,
								),
							},
						})
					} else {
						err := setRole(s, env.Get("GUILD_ID"), discord_member_id.Value.(string), env.Get("MEMBER_ROLE_ID"))
						if err != nil {
							log.Fatalf("Cannot set role: %v", err)
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								// Ignore type for now, they will be discussed in "responses"
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: fmt.Sprintf(
										msgformat,
										"Bienvenido, "+memberFound.Name,
										commandHistory,
									),
								},
							})
						}
					}

				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						// Ignore type for now, they will be discussed in "responses"
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf(
								msgformat,
								albion_name.StringValue()+" no es miembro de la guild.",
								commandHistory,
							),
						},
					})
				}

			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							msgformat,
							discord_member_id.StringValue()+", no se encontró el usuario.",
							commandHistory,
						),
					},
				})
			}

		},
	}
)

func existMemberRole(s *discordgo.Session, guildID, memberID, roleID string) bool {
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return false
	}

	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

func setRole(s *discordgo.Session, guildID, memberID, roleID string) error {
	err := s.GuildMemberRoleAdd(guildID, memberID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func Run() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot running %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, env.Get("GUILD_ID"), v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, env.Get("GUILD_ID"), v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}
