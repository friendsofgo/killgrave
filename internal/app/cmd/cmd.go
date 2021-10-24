package cmd

import (
	"github.com/friendsofgo/killgrave/internal/app/cmd/http"
	"github.com/spf13/cobra"
)

var _version = "unknown_version"

const (
	_defaultImpostersPath = "imposters"
	_defaultConfigFile    = ""
)

// NewKillgraveCmd returns cobra.Command to run killgrave command
func NewKillgraveCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "killgrave",
		Short:         "Simple way to generate mock servers",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       _version,
	}

	rootCmd.ResetFlags()
	rootCmd.PersistentFlags().StringP("imposters", "i", _defaultImpostersPath, "Directory where your imposters are located")
	rootCmd.PersistentFlags().StringP("config", "c", _defaultConfigFile, "Path to your configuration file")

	rootCmd.SetVersionTemplate("Killgrave version: {{.Version}}\n")

	rootCmd.AddCommand(http.NewHTTPCmd())

	return rootCmd
}
