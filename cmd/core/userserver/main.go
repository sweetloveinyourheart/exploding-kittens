package main

import (
	"time"

	"github.com/spf13/cobra"

	pocker_userserver "github.com/sweetloveinyourheart/planning-poker/cmd/planning-pocker/services/userserver"
	pocker_utils "github.com/sweetloveinyourheart/planning-poker/cmd/planning-pocker/utils"
	"github.com/sweetloveinyourheart/planning-poker/pkg/cmdutil"
)

const defaultShortDescription = "Planning Pocker User Server"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/planning-poker/cmd/core/userserver generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)
	commands = append(commands, pocker_userserver.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, pocker_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
