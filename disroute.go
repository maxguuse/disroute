package disroute

import (
	"errors"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	TypeCommand               = discordgo.InteractionApplicationCommand
	TypeCommandAutocompletion = discordgo.InteractionApplicationCommandAutocomplete

	TypeSubcommand      = discordgo.ApplicationCommandOptionSubCommand
	TypeSubcommandGroup = discordgo.ApplicationCommandOptionSubCommandGroup

	TypeMessageComponent = discordgo.InteractionMessageComponent
)

type DiscordCmdOption = discordgo.ApplicationCommandInteractionDataOption

type HandlerFunc func(
	*discordgo.Interaction,
	map[string]*DiscordCmdOption,
) (string, error)

type ComponentHandlerFunc func(
	*discordgo.Interaction,
) (string, error)

type Cmd struct {
	Path     string
	Handlers Handlers
	Options  []*CmdOption
}

type CmdOption struct {
	Path     string
	Handlers Handlers
	Type     discordgo.ApplicationCommandOptionType
	Options  []*CmdOption
}

type Handlers struct {
	Cmd          HandlerFunc
	Autocomplete HandlerFunc
}

type Component struct {
	Key     string
	Handler ComponentHandlerFunc
}

type Router struct {
	cmdMx sync.RWMutex
	cmds  map[string]HandlerFunc

	autocompleteMx sync.RWMutex
	autocompletes  map[string]HandlerFunc

	componentsMx  sync.RWMutex
	components    map[string]ComponentHandlerFunc
	componentFunc func(*discordgo.InteractionCreate) (key string)
}

func New(options ...func(*Router)) *Router {
	r := &Router{
		cmds:          make(map[string]HandlerFunc),
		autocompletes: make(map[string]HandlerFunc),
		components:    make(map[string]ComponentHandlerFunc),
		componentFunc: func(ic *discordgo.InteractionCreate) (key string) {
			if ic.Type != discordgo.InteractionMessageComponent {
				return ""
			}

			return ic.MessageComponentData().CustomID
		},
	}

	for _, o := range options {
		o(r)
	}

	return r
}

func WithComponentFunc(f func(*discordgo.InteractionCreate) (key string)) func(*Router) {
	return func(r *Router) {
		r.componentFunc = f
	}
}

func (r *Router) RegisterAll(cmds []*Cmd) error {
	for _, cmd := range cmds {
		pathParts := []string{cmd.Path}

		if cmd.Handlers.Cmd == nil && len(cmd.Options) == 0 {
			return errors.New("cmd has no handler and no subcommands")
		}

		if len(cmd.Options) == 0 {
			r.registerCmd(pathParts, cmd.Handlers)
		}

		for _, opt := range cmd.Options {
			pathParts = append(pathParts, opt.Path)

			if opt.Type == TypeSubcommand {
				if opt.Handlers.Cmd == nil {
					return errors.New("subcommand has no handler")
				}

				r.registerCmd(pathParts, opt.Handlers)
			}

			if opt.Type == TypeSubcommandGroup {
				for _, sub := range opt.Options {
					if sub.Type != TypeSubcommand {
						continue
					}

					if sub.Handlers.Cmd == nil {
						return errors.New("subcommand has no handler")
					}

					pathParts = append(pathParts, sub.Path)

					r.registerCmd(pathParts, sub.Handlers)

					pathParts = pathParts[:len(pathParts)-1]
				}
			}

			pathParts = pathParts[:len(pathParts)-1]
		}
	}

	return nil
}

func (r *Router) registerCmd(pathParts []string, hs Handlers) {
	path := strings.Join(pathParts, ":")

	r.cmdMx.Lock()
	defer r.cmdMx.Unlock()

	r.cmds[path] = hs.Cmd

	if hs.Autocomplete != nil {
		r.autocompleteMx.Lock()
		defer r.autocompleteMx.Unlock()

		r.autocompletes[path] = hs.Autocomplete
	}
}

func (r *Router) GetAll() map[string]HandlerFunc {
	r.cmdMx.RLock()
	defer r.cmdMx.RUnlock()

	return r.cmds
}

func (r *Router) FindAndExecute(i *discordgo.InteractionCreate) (string, error) {
	if i.Type != TypeCommand {
		return "", errors.New("invalid interaction type")
	}

	r.cmdMx.RLock()
	defer r.cmdMx.RUnlock()

	hd := r.buildHandlerData(i)

	if h, ok := r.cmds[hd.path]; ok {
		return h(i.Interaction, hd.opts)
	}

	return "", errors.New("command not registered")
}

func (r *Router) FindAndAutocomplete(i *discordgo.InteractionCreate) (string, error) {
	if i.Type != TypeCommandAutocompletion {
		return "", errors.New("invalid interaction type")
	}

	r.autocompleteMx.RLock()
	defer r.autocompleteMx.RUnlock()

	hd := r.buildHandlerData(i)

	if h, ok := r.autocompletes[hd.path]; ok {
		return h(i.Interaction, hd.opts)
	}

	return "", errors.New("autocompletion not registered")
}

func (r *Router) RegisterComponents(cmps []*Component) error {
	r.componentsMx.Lock()
	defer r.componentsMx.Unlock()

	for _, cmp := range cmps {
		if strings.TrimSpace(cmp.Key) == "" || cmp.Handler == nil {
			return errors.New("invalid component, missing key or handler")
		}
		r.components[cmp.Key] = cmp.Handler
	}

	return nil
}

func (r *Router) FindComponentAndExecute(i *discordgo.InteractionCreate) (string, error) {
	if i.Type != TypeMessageComponent {
		return "", errors.New("invalid interaction type")
	}

	r.componentsMx.RLock()
	defer r.componentsMx.RUnlock()

	key := r.componentFunc(i)

	if h, ok := r.components[key]; ok {
		return h(i.Interaction)
	}

	return "", errors.New("component not registered")
}

type handlerData struct {
	path string
	opts map[string]*DiscordCmdOption
}

func (r *Router) buildHandlerData(i *discordgo.InteractionCreate) *handlerData {
	d := i.ApplicationCommandData()

	pathParts := []string{d.Name}
	options := r.buildOptionsMap(d.Options)

	if len(d.Options) == 0 {
		return &handlerData{
			path: strings.Join(pathParts, ":"),
			opts: options,
		}
	}

	if d.Options[0].Type == TypeSubcommand {
		pathParts = append(pathParts, d.Options[0].Name)
		options = r.buildOptionsMap(d.Options[0].Options)
	}

	if d.Options[0].Type == TypeSubcommandGroup {
		pathParts = append(pathParts,
			d.Options[0].Name,
			d.Options[0].Options[0].Name,
		)
		options = r.buildOptionsMap(d.Options[0].Options[0].Options)
	}

	return &handlerData{
		path: strings.Join(pathParts, ":"),
		opts: options,
	}
}

func (r *Router) buildOptionsMap(options []*DiscordCmdOption) map[string]*DiscordCmdOption {
	commandOptions := make(map[string]*DiscordCmdOption)
	for _, option := range options {
		commandOptions[option.Name] = option
	}

	return commandOptions
}
