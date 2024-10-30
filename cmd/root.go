/*
Copyright © 2024 Saman Dehestani <github.com/drippypale>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-sharif-net",
	Short: "A cli tool to make the net2 authentication easier",
	Long:  `A cli tool to make the net2 authentication easier`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-sharif-net.yaml)")
	rootCmd.PersistentFlags().Bool("use-ip", false, "if specified, uses the provided IP in the config instead of the domain.")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".go-sharif-net" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-sharif-net")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		setDefaultConfig()
	}
}

func setDefaultConfig() {
	viper.Set("hostIP", "https://172.17.1.214")
	viper.Set("hostDomain", "https://net2.sharif.edu")
	viper.Set("loginEndpoint", "/login")
	viper.Set("statusEndpoint", "/status")
	viper.Set("logoutEndpoint", "/logout")

	err2 := viper.SafeWriteConfig()
	cobra.CheckErr(err2)
}
