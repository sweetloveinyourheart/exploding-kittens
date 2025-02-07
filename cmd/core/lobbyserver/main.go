package main

import (
	"time"

	"github.com/spf13/cobra"

	kittens_lobbyserver "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/lobbyserver"

	kittens_utils "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/utils"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
)

const defaultShortDescription = "Exploding Kittens Lobby Server"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/exploding-kittens/cmd/core/lobbyserver generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)
	commands = append(commands, kittens_lobbyserver.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, kittens_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
