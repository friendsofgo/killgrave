package app

import (
	"github.com/friendsofgo/killgrave/internal/app/cmd"
)

// Run creates and executes new killgrave command
func Run() error {
	rootCmd := cmd.NewKillgraveCmd()
	return rootCmd.Execute()
}
