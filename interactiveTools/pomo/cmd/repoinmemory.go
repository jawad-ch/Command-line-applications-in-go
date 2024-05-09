package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jawad-ch/Command-line-application-in-go/interactiveTools/pomo/pomodoro"
	"github.com/jawad-ch/Command-line-application-in-go/interactiveTools/pomo/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// repoinmemoryCmd represents the repoinmemory command
var repoinmemoryCmd = &cobra.Command{
	Use:   "repoinmemory",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repoinmemory called")
	},
}
var cfgFile string

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".pScan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pomo")
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	rootCmd.AddCommand(repoinmemoryCmd)

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.pomo.yaml)")
	rootCmd.Flags().DurationP("pomo", "p", 25*time.Minute,
		"Pomodoro duration")
	rootCmd.Flags().DurationP("short", "s", 5*time.Minute,
		"Short break duration")
	rootCmd.Flags().DurationP("long", "l", 15*time.Minute,
		"Long break duration")
	viper.BindPFlag("pomo", rootCmd.Flags().Lookup("pomo"))
	viper.BindPFlag("short", rootCmd.Flags().Lookup("short"))
	viper.BindPFlag("long", rootCmd.Flags().Lookup("long"))
}

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
