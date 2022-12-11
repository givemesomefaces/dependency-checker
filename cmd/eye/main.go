package main

import (
	"os"

	"eyes/commands"
	"eyes/internal/logger"
)

func main() {
	if err := commands.Execute(); err != nil {
		logger.Log.Errorln(err)
		os.Exit(1)
	}
}
