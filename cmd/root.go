/*
Copyright © 2025 Motalleb Fallahnezhad

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/fmotalleb/go-tools/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/controller"
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/system"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "the-one",
	Short: "A simple init system for monolithic containers",
	Long: `Simple yet fast init system for monolithic containers.
It is designed to be lightweight and easy to use, making it ideal for
containers that require a simple init system.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		isVerbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		var l *zap.Logger
		if l, err = buildLogger(isVerbose); err != nil {
			return err
		}
		logging.SetRootLogger(l)
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		isVerbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		cfg := new(config.Config)
		if err := config.Parse(cfg, cfgFile, isVerbose); err != nil {
			return err
		}
		return controller.Boot(system.NewSystemContext(), cfg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(
		"config",
		"c",
		"config.toml",
		"config file (default is ./config.toml)",
	)

	rootCmd.PersistentFlags().BoolP(
		"verbose",
		"v",
		false,
		"enable verbose development logger instead of JSON",
	)
}

func buildLogger(isDebug bool) (*zap.Logger, error) {
	if isDebug {
		log.SetDebugDefaults()
	}
	b := log.
		NewBuilder().
		FromEnv()
	return b.Build()
}
