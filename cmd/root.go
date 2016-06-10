package cmd

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var cfgFile string
var doToken string
var doDomain string
var checkUrlList []string

var displayVersion string
var showVersion bool
var verbose bool
var debug bool

// RootCmd handles all of our arguments/options
var RootCmd = &cobra.Command{
	Use:   "digitalocean-ddns",
	Short: "Use digitalocean as your DDNS service",
	Long: `A service application that watches your external IP for changes
and updates a DigitalOcean domain when a change is detected`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(displayVersion)
			os.Exit(1)
		}

		log.SetFormatter(&prefixed.TextFormatter{
			TimestampFormat: time.RFC3339,
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.InfoLevel)

		if debug {
			log.SetLevel(log.DebugLevel)
		}
		log.Infof("config: file=%s", viper.ConfigFileUsed())
		log.Debugf("config: domain=%s token=%s...",
			viper.GetString("domain"),
			viper.GetString("token")[:5])

		RunSync(viper.GetString("token"), viper.GetString("domain"))
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
	RootCmd.PersistentFlags().StringVarP(&doToken, "token", "t", "",
		"the DigitalOcean username to authenticate with")
	RootCmd.PersistentFlags().StringVarP(&doDomain, "domain", "d", "",
		"the DigitalOcean domain to update")
	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version number and exit.")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false,
		"Print logs to stdout instead of file")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"Include debug statements in log output")

	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("domain", RootCmd.PersistentFlags().Lookup("domain"))

	viper.SetDefault("url_list", []string{
		"http://icanhazip.com",
		"http://whatismyip.akamai.com/",
		"http://whatsmyip.me/",
		"http://wtfismyip.com/text",
		"http://api.ipify.org/",
		"http://ip.catnapgames.com",
		"http://ip.ryansanden.com",
	})
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
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error opening config: ", err)
	}
}
