package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

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
		{
			Name:        "crear-evento",
			Description: "Crear un evento",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"registrar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.ChannelID != env.Get("COMMANDS_CHANNEL_ID") {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "No puedes usar este comando en este canal.",
					},
				})

				time.Sleep(time.Second * 10)

				s.InteractionResponseDelete(i.Interaction)
			}
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
					setNewMemberNickname(s, env.Get("GUILD_ID"), discord_member_id.Value.(string), albion_name.StringValue())
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
		"crear-evento": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "create_event_" + i.Interaction.Member.User.ID,
					Title:    "Crear evento",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "title",
									Label:       "Titulo",
									Style:       discordgo.TextInputShort,
									Placeholder: "Ingrese el titulo del evento",
									Required:    true,
									MaxLength:   300,
									MinLength:   5,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "description",
									Label:       "Descripción",
									Style:       discordgo.TextInputParagraph,
									Placeholder: "Ingrese la descripción del evento",
									Required:    true,
									MaxLength:   2000,
									MinLength:   10,
								}},
						}},
				},
			})

			if err != nil {
				log.Fatalf("Error al crear el evento: %v", err)
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

func setNewMemberNickname(s *discordgo.Session, guildID, memberID, nickname string) error {
	log.Println(nickname)
	err := s.GuildMemberNickname(guildID, memberID, nickname)
	if err != nil {
		return err
	}
	return nil
}

func createGuildEvent(s *discordgo.Session, data discordgo.ModalSubmitInteractionData) error {
	title := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	description := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	startingTime := time.Now().Add(1 * time.Hour)

	endingTime := time.Now().Add(1 * time.Hour)

	_, err := s.GuildScheduledEventCreate(env.Get("GUILD_ID"), &discordgo.GuildScheduledEventParams{
		Name:               title,
		Description:        description,
		ScheduledStartTime: &startingTime,
		ScheduledEndTime:   &endingTime,
		EntityType:         discordgo.GuildScheduledEventEntityTypeVoice,
		ChannelID:          env.Get("EVENT_CHANNEL_VOICE_ID"),
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	})
	if err != nil {
		log.Printf("Error creating scheduled event: %v", err)
		return err
	}
	return nil
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "create_event") {
				return
			}

			err := createGuildEvent(s, data)
			if err != nil {
				panic(err)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Se ha creado el evento.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
			userid := strings.Split(data.CustomID, "_")[2]
			_, err = s.ChannelMessageSend(env.Get("COMMANDS_CHANNEL_ID"), fmt.Sprintf(
				"Evento creado por <@%s>\n\n**Evento**:\n%s\n\n**Descripción**:\n%s",
				userid,
				data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			))
			if err != nil {
				panic(err)
			}

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
