/*
Copyright © 2020 zhijiezhang

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hitzhangjie/log"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gorpc",
	Short: config.LoadTranslation("rootCmdUsage", nil),
	Long:  config.LoadTranslation("rootCmdUsageLong", nil),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "apiLoadConfig", "", "apiLoadConfig file (default is $HOME/.hitzhangjie-cmdline.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in apiLoadConfig file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use apiLoadConfig file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search apiLoadConfig in home directory with name ".hitzhangjie-cmdline" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hitzhangjie-cmdline")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a apiLoadConfig file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using apiLoadConfig file:", viper.ConfigFileUsed())
	}
}

func withErrorCheck(stage string, err error) {
	if err != nil {
		log.Error(stage+" error: %v", err)
		os.Exit(1)
	}
}
