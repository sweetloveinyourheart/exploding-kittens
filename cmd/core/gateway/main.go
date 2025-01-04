package main

import (
	"time"

	"github.com/spf13/cobra"

	kittens_gateway "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/services/gateway"
	kittens_utils "github.com/sweetloveinyourheart/exploding-kittens/cmd/exploding-kittens/utils"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
)

const defaultShortDescription = "EXPLODING KITTENS API Gateway"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/exploding-kittens/cmd/core/gateway generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)
	commands = append(commands, kittens_gateway.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, kittens_utils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
