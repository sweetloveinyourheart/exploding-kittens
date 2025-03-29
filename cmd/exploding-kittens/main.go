package main

import (
	"time"

	"github.com/spf13/cobra"

	kittens_clientserver "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/clientserver"
	kittens_gameengineserver "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/gameengineserver"
	kittens_lobbyserver "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/lobbyserver"
	kittens_userserver "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/userserver"
	kittens_utils "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/utils"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
)

//go:generate go run github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens generate

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

func main() {
	commands := make([]*cobra.Command, 0)

	commands = append(commands, kittens_clientserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, kittens_lobbyserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, kittens_userserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, kittens_gameengineserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, kittens_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
