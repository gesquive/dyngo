package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var doUserName string
var doPassword string
var doDomain string
var displayVersion string
var showVersion bool
var verbose bool
var debug bool

var RootCmd = &cobra.Command{
	Use:   "digitalocean-ddns",
	Short: "Use digitalocean as your DDNS service",
	Long: `A service application that watches your external IP for changes
and updates a DigitalOcean domain when a change is detected`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(displayVersion)
			os.Exit(1)
		}
		// fmt.Println(viper.Get("username"))
		fmt.Printf("Username: %s\n", viper.Get("username"))
		fmt.Printf("Password: %s\n", viper.Get("password"))
		fmt.Printf("Domain  : %s\n", viper.Get("domain"))
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is ./config.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().StringVarP(&doUserName, "username", "u", "",
		"the DigitalOcean username to authenticate with")
	RootCmd.PersistentFlags().StringVarP(&doPassword, "password", "p", "",
		"the DigitalOcean password to authenticate with")
	RootCmd.PersistentFlags().StringVarP(&doDomain, "domain", "d", "",
		"the DigitalOcean domain to update")
	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version number and exit.")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false,
		"Print logs to stdout")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"Print out debug statements")

	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("domain", RootCmd.PersistentFlags().Lookup("domain"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config.yml") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/digitalocean-ddns") // adding home directory as first search path
	viper.AddConfigPath("/etc/digitalocean-ddns")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
