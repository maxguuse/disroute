package disroute

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	TypeSubcommand      = discordgo.ApplicationCommandOptionSubCommand
	TypeSubcommandGroup = discordgo.ApplicationCommandOptionSubCommandGroup
)

type DiscordCmdOption = discordgo.ApplicationCommandInteractionDataOption

type HandlerFunc func(
	*discordgo.Interaction,
	map[string]*DiscordCmdOption,
) error

type Cmd struct {
	Path    string
	Handler HandlerFunc
	Options []*CmdOption
}

type CmdOption struct {
	Path    string
	Handler HandlerFunc
	Type    discordgo.ApplicationCommandOptionType
	Options []*CmdOption
}

type Router struct {
	cmds map[string]HandlerFunc
}

func New() *Router {
	return &Router{
		cmds: make(map[string]HandlerFunc),
	}
}

func (r *Router) RegisterAll(cmds []*Cmd) error {
	if len(cmds) == 0 {
		return nil
	}

	for _, cmd := range cmds {
		pathParts := []string{cmd.Path}

		if cmd.Handler == nil && len(cmd.Options) == 0 {
			return errors.New("cmd has no handler and no subcommands")
		}

		if len(cmd.Options) == 0 {
			r.cmds[cmd.Path] = cmd.Handler
			continue
		}

		for _, opt := range cmd.Options {
			if opt.Type != TypeSubcommand && opt.Type != TypeSubcommandGroup {
				return errors.New("cmd has non-subcommand option")
			}

			pathParts = append(pathParts, opt.Path)

			if opt.Type == TypeSubcommand {
				if opt.Handler == nil {
					return errors.New("subcommand has no handler")
				}

				p := strings.Join(pathParts, ":")
				h := opt.Handler

				r.cmds[p] = h

				pathParts = pathParts[:len(pathParts)-1]

				continue
			}

			if opt.Type == TypeSubcommandGroup {
				if len(opt.Options) == 0 {
					return errors.New("subcommand group has no subcommands")
				}

				for _, sub := range opt.Options {
					if sub.Type != TypeSubcommand {
						return errors.New("subcommand group has non-subcommand option")
					}

					if sub.Handler == nil {
						return errors.New("subcommand has no handler")
					}

					pathParts = append(pathParts, sub.Path)

					p := strings.Join(pathParts, ":")
					h := sub.Handler

					r.cmds[p] = h

					pathParts = pathParts[:len(pathParts)-1]
				}

				pathParts = pathParts[:len(pathParts)-1]
			}
		}
	}

	return nil
}

func (r *Router) GetAll() map[string]HandlerFunc {
	return r.cmds

}

func (r *Router) FindAndExecute(i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionApplicationCommand {
		return errors.New("invalid interaction type")
	}

	data := i.ApplicationCommandData()

	pathParts := []string{data.Name}
	options := r.buildOptionsMap(data.Options)

	for _, opt := range data.Options {
		if opt.Type == TypeSubcommand {
			pathParts = append(pathParts, opt.Name)
			options = r.buildOptionsMap(opt.Options)
			break
		}

		if opt.Type == TypeSubcommandGroup {
			pathParts = append(pathParts, opt.Name)
			for _, subOpt := range opt.Options {
				if subOpt.Type == TypeSubcommand {
					pathParts = append(pathParts, subOpt.Name)
					options = r.buildOptionsMap(subOpt.Options)
					break
				}
			}
		}
	}

	path := strings.Join(pathParts, ":")

	var h HandlerFunc
	var ok bool
	if h, ok = r.cmds[path]; !ok {
		return errors.New("command not registered")
	}

	return h(i.Interaction, options)
}

func (r *Router) buildOptionsMap(options []*DiscordCmdOption) map[string]*DiscordCmdOption {
	commandOptions := make(map[string]*DiscordCmdOption)
	for _, option := range options {
		commandOptions[option.Name] = option
	}

	return commandOptions
}
