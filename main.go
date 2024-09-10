package main

import (
	"os"

	"github.com/gptscript-ai/cmd"
	"github.com/gptscript-ai/otto/pkg/cli"
)

func main() {
	// Don't shutdown on SIGTERM, only on SIGINT. SIGTERM is handled by the controller leader election
	cmd.ShutdownSignals = []os.Signal{os.Interrupt}
	cmd.Main(cli.New())
}