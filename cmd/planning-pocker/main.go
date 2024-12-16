package main

import (
	"time"

	"github.com/spf13/cobra"

	pocker_userserver "github.com/sweetloveinyourheart/planning-pocker/cmd/planning-pocker/services/userserver"
	pocker_utils "github.com/sweetloveinyourheart/planning-pocker/cmd/planning-pocker/utils"
	"github.com/sweetloveinyourheart/planning-pocker/pkg/cmdutil"
)

//go:generate go run github.com/sweetloveinyourheart/planning-pocker/cmd/planning-pocker generate

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

func main() {
	commands := make([]*cobra.Command, 0)

	commands = append(commands, pocker_userserver.Command(cmdutil.ServiceRootCmd))
	commands = append(commands, pocker_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
