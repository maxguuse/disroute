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

	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd: EmptyHandler,
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if len(r.GetAll()) != 1 {
		t.Errorf("Expected 1 command, got %d", len(r.GetAll()))
	}
}

func Test_RegisterAll_Subcommands(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "sub",
			Options: []*disroute.CmdOption{
				{
					Path: "cmd",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd: EmptyHandler,
					},
				},
				{
					Path: "cmd2",
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

	if len(r.GetAll()) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(r.GetAll()))
	}
}

func Test_RegisterAll_SubcommandGroups(t *testing.T) {
	r := disroute.New()

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
								Cmd: EmptyHandler,
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd: EmptyHandler,
							},
						},
					},
				},
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Error(err)
	}
}

func Test_RegisterAll_MixedSubcommands(t *testing.T) {
	r := disroute.New()

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
								Cmd: EmptyHandler,
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd: EmptyHandler,
							},
						},
					},
				},
				{
					Path: "cmd",
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

	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd: EmptyHandler,
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	_, err = r.FindAndExecute(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd",
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
	}
}

func Test_FindAndExecute_Subcommand(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
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
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd3",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name: "sub",
						Type: discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
	}
}
func Test_FindAndExecute_SubcommandGroup(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "cmd4",
			Options: []*disroute.CmdOption{
				{
					Path: "gr",
					Type: disroute.TypeSubcommandGroup,
					Options: []*disroute.CmdOption{
						{
							Path: "sub2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd: EmptyHandler,
							},
						},
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
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd4",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name: "gr",
						Type: discordgo.ApplicationCommandOptionSubCommandGroup,
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Name: "sub2",
								Type: discordgo.ApplicationCommandOptionSubCommand,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
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

	cmds := []*disroute.Cmd{
		{
			Path: "cmd",
			Handlers: disroute.Handlers{
				Cmd:          EmptyHandler,
				Autocomplete: EmptyHandler,
			},
		},
	}

	err := r.RegisterAll(cmds)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	_, err = r.FindAndAutocomplete(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd",
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
	}
}

func Test_FindAndAutocomplete_Subcommand(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "cmd3",
			Options: []*disroute.CmdOption{
				{
					Path: "sub",
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          EmptyHandler,
						Autocomplete: EmptyHandler,
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
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd3",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name: "sub",
						Type: discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
	}
}
func Test_FindAndAutocomplete_SubcommandGroup(t *testing.T) {
	r := disroute.New()

	cmds := []*disroute.Cmd{
		{
			Path: "cmd4",
			Options: []*disroute.CmdOption{
				{
					Path: "gr",
					Type: disroute.TypeSubcommandGroup,
					Options: []*disroute.CmdOption{
						{
							Path: "sub2",
							Type: disroute.TypeSubcommand,
							Handlers: disroute.Handlers{
								Cmd:          EmptyHandler,
								Autocomplete: EmptyHandler,
							},
						},
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
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "cmd4",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name: "gr",
						Type: discordgo.ApplicationCommandOptionSubCommandGroup,
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Name: "sub2",
								Type: discordgo.ApplicationCommandOptionSubCommand,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Error("Expected nil error, got", err)
	}
}
