package main

import (
	"time"

	"github.com/spf13/cobra"

	pp_userserver "github.com/sweetloveinyourheart/planning-poker/cmd/planning-pocker/services/userserver"
	pp_utils "github.com/sweetloveinyourheart/planning-poker/cmd/planning-pocker/utils"
	"github.com/sweetloveinyourheart/planning-poker/pkg/cmdutil"
)

//go:generate go run github.com/sweetloveinyourheart/planning-poker/cmd/planning-pocker generate

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

func main() {
	commands := make([]*cobra.Command, 0)

	commands = append(commands, pp_userserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, pp_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
