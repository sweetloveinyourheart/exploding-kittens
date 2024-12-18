package main

import (
	"time"

	"github.com/spf13/cobra"

	pocker_gateway "github.com/sweetloveinyourheart/planning-pocker/cmd/planning-pocker/services/gateway"
	pocker_utils "github.com/sweetloveinyourheart/planning-pocker/cmd/planning-pocker/utils"
	"github.com/sweetloveinyourheart/planning-pocker/pkg/cmdutil"
)

const defaultShortDescription = "Planning Pocker API Gateway"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/planning-pocker/cmd/core/gateway generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)
	commands = append(commands, pocker_gateway.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, pocker_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
