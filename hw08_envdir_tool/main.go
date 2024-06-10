package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-envdir [path to env dir] [command]",
	Short: "Утилита позволяет запускать программы, получая переменные окружения из определенной директории",
	Args:  cobra.MinimumNArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		dir := args[0]
		command := args[1:]

		env, err := ReadDir(dir)
		if err != nil {
			slog.Error("failed reading directory", "error", err)
			os.Exit(1)
		}

		exitCode := RunCmd(command, env)
		os.Exit(exitCode)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed executing command", "error", err)
	}
}
