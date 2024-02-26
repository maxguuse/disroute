package disroute_test

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/disroute"
)

var (
	EmptyHandler = func(
		*discordgo.Interaction,
		map[string]*disroute.DiscordCmdOption,
	) (string, error) {
		return "", nil
	}
	ErrorHandler = func(
		*discordgo.Interaction,
		map[string]*disroute.DiscordCmdOption,
	) (string, error) {
		return "", errors.New("error")
	}
)

type table struct {
	cmds         []*disroute.Cmd
	interactions []*discordgo.InteractionCreate
}

func getSingleCmd(iType discordgo.InteractionType) table {
	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd:          EmptyHandler,
				Autocomplete: EmptyHandler,
			},
		},
	}

	i := []*discordgo.InteractionCreate{
		{
			Interaction: &discordgo.Interaction{
				Type: iType,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "cmd",
				},
			},
		},
	}

	return table{
		cmds:         cmds,
		interactions: i,
	}
}

func getSubcommandCmd(iType discordgo.InteractionType) table {
	cmds := []*disroute.Cmd{
		{
			Path: "sub",
			Options: []*disroute.CmdOption{
				{
					Path: "cmd",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          EmptyHandler,
						Autocomplete: EmptyHandler,
					},
				},
				{
					Path: "cmd2",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          EmptyHandler,
						Autocomplete: EmptyHandler,
					},
				},
			},
		},
	}
	i := []*discordgo.InteractionCreate{
		{
			Interaction: &discordgo.Interaction{
				Type: iType,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "sub",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name: "cmd",
							Type: discordgo.ApplicationCommandOptionSubCommand,
						},
					},
				},
			},
		},
	}

	return table{
		cmds:         cmds,
		interactions: i,
	}
}

func getSubcommandGroupCmd(iType discordgo.InteractionType) table {
	cmds := []*disroute.Cmd{
		{
			Path: "gr-sub",
			Options: []*disroute.CmdOption{
				{
					Path: "gr",
					Type: disroute.TypeSubcommandGroup,
					Options: []*disroute.CmdOption{
						{
							Path: "cmd",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd:          EmptyHandler,
								Autocomplete: EmptyHandler,
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd:          EmptyHandler,
								Autocomplete: EmptyHandler,
							},
						},
						{
							Path: "non_subcommand",
							Type: discordgo.ApplicationCommandOptionBoolean,
						},
					},
				},
			},
		},
	}

	i := []*discordgo.InteractionCreate{
		{
			Interaction: &discordgo.Interaction{
				Type: iType,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "gr-sub",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name: "gr",
							Type: discordgo.ApplicationCommandOptionSubCommandGroup,
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{
									Name: "cmd",
									Type: discordgo.ApplicationCommandOptionSubCommand,
								},
							},
						},
					},
				},
			},
		},
	}

	return table{
		cmds:         cmds,
		interactions: i,
	}
}

func getMixedCmd(iType discordgo.InteractionType) table {
	cmds := []*disroute.Cmd{
		{
			Path: "gr-sub",
			Options: []*disroute.CmdOption{
				{
					Path: "gr",
					Type: disroute.TypeSubcommandGroup,
					Options: []*disroute.CmdOption{
						{
							Path: "cmd",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd:          EmptyHandler,
								Autocomplete: EmptyHandler,
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd:          EmptyHandler,
								Autocomplete: EmptyHandler,
							},
						},
					},
				},
				{
					Path: "cmd",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          EmptyHandler,
						Autocomplete: EmptyHandler,
					},
				},
			},
		},
	}

	i := []*discordgo.InteractionCreate{
		{
			Interaction: &discordgo.Interaction{
				Type: iType,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "gr-sub",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name: "cmd",
							Type: discordgo.ApplicationCommandOptionSubCommand,
						},
					},
				},
			},
		},
		{
			Interaction: &discordgo.Interaction{
				Type: iType,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "gr-sub",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name: "gr",
							Type: discordgo.ApplicationCommandOptionSubCommandGroup,
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{
									Name: "cmd",
									Type: discordgo.ApplicationCommandOptionSubCommand,
								},
							},
						},
					},
				},
			},
		},
	}

	return table{
		cmds:         cmds,
		interactions: i,
	}
}

func Test_RegisterAll_Errors(t *testing.T) {
	r := disroute.New()

	// Testing for the case when cmds is empty
	err := r.RegisterAll([]*disroute.Cmd{})
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Testing for the case when a cmd has no handler and no subcommands
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test"}})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Testing for the case when a subcommand has no handler
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommand}}}})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Testing for the case when a subcommand group has a subcommand with nil handler
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommandGroup, Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommand}}}}}})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func Test_RegisterAll_SingleCmd(t *testing.T) {
	r := disroute.New()

	cases := getSingleCmd(0)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(r.GetAll()) != 1 {
		t.Errorf("Expected 1 command, got %d", len(r.GetAll()))
	}
}

func Test_RegisterAll_Subcommands(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandCmd(0)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(r.GetAll()) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(r.GetAll()))
	}
}

func Test_RegisterAll_SubcommandGroups(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandGroupCmd(0)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Error(err)
	}
}

func Test_RegisterAll_MixedSubcommands(t *testing.T) {
	r := disroute.New()

	cases := getMixedCmd(0)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Error(err)
	}

	if len(r.GetAll()) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(r.GetAll()))
	}
}

func Test_FindAndExecute_Errors(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd: ErrorHandler,
			},
		},
		{
			Path: "cmd3",
			Options: []*disroute.CmdOption{
				{
					Path: "sub",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd: EmptyHandler,
					},
				},
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	_, err = r.FindAndExecute(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndExecute(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndExecute(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd2",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndExecute(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd3",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func Test_FindAndExecute_SingleCmd(t *testing.T) {
	r := disroute.New()

	cases := getSingleCmd(disroute.TypeCommand)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndExecute(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}

func Test_FindAndExecute_Subcommand(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandCmd(disroute.TypeCommand)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndExecute(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}
func Test_FindAndExecute_SubcommandGroup(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandGroupCmd(disroute.TypeCommand)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndExecute(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}
func Test_FindAndAutocomplete_Errors(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd: ErrorHandler,
			},
		},
		{
			Path: "cmd3",
			Options: []*disroute.CmdOption{
				{
					Path: "sub",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd: EmptyHandler,
					},
				},
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	_, err = r.FindAndAutocomplete(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndAutocomplete(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndAutocomplete(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd2",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = r.FindAndAutocomplete(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd3",
			},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func Test_FindAndAutocomplete_SingleCmd(t *testing.T) {
	r := disroute.New()

	cases := getSingleCmd(disroute.TypeCommandAutocompletion)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndAutocomplete(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}

func Test_FindAndAutocomplete_Subcommand(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandCmd(disroute.TypeCommandAutocompletion)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndAutocomplete(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}
func Test_FindAndAutocomplete_SubcommandGroup(t *testing.T) {
	r := disroute.New()

	cases := getSubcommandGroupCmd(disroute.TypeCommandAutocompletion)

	err := r.RegisterAll(cases.cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	for _, i := range cases.interactions {
		_, err = r.FindAndAutocomplete(i)
		if err != nil {
			t.Error("Expected nil error, got", err)
		}
	}
}
