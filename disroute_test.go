package disroute_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/disroute"
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

	// Testing for the case when a cmd has a non-subcommand option
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: discordgo.ApplicationCommandOptionInteger}}}})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Testing for the case when a subcommand has no handler
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommand}}}})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Testing for the case when a subcommand group has no subcommands
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommandGroup}}}})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Testing for the case when a subcommand group has a non-subcommand option
	err = r.RegisterAll([]*disroute.Cmd{{Path: "test", Options: []*disroute.CmdOption{{Type: disroute.TypeSubcommandGroup, Options: []*disroute.CmdOption{{Type: discordgo.ApplicationCommandOptionInteger}}}}}})
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
			Handler: func(
				*discordgo.Interaction,
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
			) {
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
					Handler: func(
						*discordgo.Interaction,
						map[string]*discordgo.ApplicationCommandInteractionDataOption,
					) {
					},
				},
				{
					Path: "cmd2",
					Type: disroute.TypeSubcommand,
					Handler: func(
						*discordgo.Interaction,
						map[string]*discordgo.ApplicationCommandInteractionDataOption,
					) {
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
							Handler: func(
								*discordgo.Interaction,
								map[string]*discordgo.ApplicationCommandInteractionDataOption,
							) {
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handler: func(
								*discordgo.Interaction,
								map[string]*discordgo.ApplicationCommandInteractionDataOption,
							) {
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
							Handler: func(
								*discordgo.Interaction,
								map[string]*discordgo.ApplicationCommandInteractionDataOption,
							) {
							},
						},
						{
							Path: "cmd2",
							Type: disroute.TypeSubcommand,
							Handler: func(
								*discordgo.Interaction,
								map[string]*discordgo.ApplicationCommandInteractionDataOption,
							) {
							},
						},
					},
				},
				{
					Path: "cmd",
					Type: disroute.TypeSubcommand,
					Handler: func(
						*discordgo.Interaction,
						map[string]*discordgo.ApplicationCommandInteractionDataOption,
					) {
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
