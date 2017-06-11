// Copyright © 2017 Yehor Nazarkin <nimnull@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"node_agent/agent"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "node_agent",
	Short: "Node health and usage report agent",
	Long: `	Agent that reports node name and number of connections to the
	main API server`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRun: func(cmd *cobra.Command, args []string) {
		apiURI := agent.ApiURIBuilder(
			viper.GetString("api"),
			viper.GetBool("ssl"))
		viper.Set("api", apiURI)
	},
	Run: func(cmd *cobra.Command, args []string) {
		agent.StartReactor()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is node_agent_conf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().Bool("ssl", false, "Switch SSL usage")
	RootCmd.Flags().Bool("debug", false, "Switch DEBUG mode")
	RootCmd.Flags().String("api", "localhost", "FQDN for api to report")
	RootCmd.Flags().String("nodename", hostname, "FQDN to report as current node")
	RootCmd.Flags().Uint("port", 3000, "Port to observe")
	RootCmd.Flags().Uint("rtime", 3, "Refresh timeout seconds")

	viper.BindPFlag("ssl", RootCmd.Flags().Lookup("ssl"))
	viper.BindPFlag("debug", RootCmd.Flags().Lookup("debug"))
	viper.BindPFlag("api", RootCmd.Flags().Lookup("api"))
	viper.BindPFlag("nodename", RootCmd.Flags().Lookup("nodename"))
	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("rtime", RootCmd.Flags().Lookup("rtime"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "node_agent_conf" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("node_agent_conf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})
}
