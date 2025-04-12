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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/fmotalleb/the-one/config"
)

var (
	cfgFile string
	cfg     config.Config
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
	Run: func(cmd *cobra.Command, args []string) {
		println("the-one called")
		data, err := yaml.Marshal(cfg)
		fmt.Printf("%s\n%v\n%v", data, err, cfg.Services["simple"].Lazy)
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.seed.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	c, err := config.ReadAndMergeConfig(cfgFile)
	if err != nil {
		panic(err)
	}
	cfg, err = config.DecodeConfig(c)
	if err != nil {
		panic(err)
	}
}
