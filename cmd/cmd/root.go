package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cdb",
	Short: "work with calibre's db",
}

func debug(cmd *cobra.Command, args []string) {
	if cmd.Flags().Changed("debug") {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cdb/config.toml)")

	rootCmd.PersistentFlags().StringP("lib", "l", "", "library by name")
	viper.BindPFlag("lib", rootCmd.PersistentFlags().Lookup("lib"))

	rootCmd.PersistentFlags().Bool("dry-run", false, "show generated command")
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().Bool("debug", false, "debug commands")
	rootCmd.Flags().MarkHidden("debug")

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

		cfgDir := filepath.Join(home, ".config/cdb/")
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.SetDefault("lib", viper.GetString("options.default"))
	} else {
		log.Fatalf("config error %v\n", err)
	}
}
