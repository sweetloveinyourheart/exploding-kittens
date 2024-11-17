package main

import (
	"time"

	"github.com/spf13/cobra"
	pp_userserver "github.com/sweetloveinyourheart/planning-poker/cmd/planing-pocker/services/userserver"
	pp_utils "github.com/sweetloveinyourheart/planning-poker/cmd/planing-pocker/utils"
	"github.com/sweetloveinyourheart/planning-poker/pkg/cmdutil"
)

const defaultShortDescription = "Planning Pocker User Server"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/planning-poker/core/userserver generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)
	commands = append(commands, pp_userserver.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, pp_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
