/*
Copyright © 2024 superryanguo
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/superryanguo/ryai/utils"
)

var (
	cfgFile string
	logger  *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "ryai",
	Short: "A ryai app",
	Long:  `A ryai application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ryai cfgFile=%s\n", cfgFile)
		logger.Info("Application started", slog.String("ryai cfgFile", cfgFile))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./conf/ryai.yaml", "config file (default is ./conf/ryai.yaml)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number")

	rootCmd.AddCommand(versionCmd)

	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s, Built@%s\n", utils.Version, utils.BuildTime)
	},
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if v, _ := rootCmd.PersistentFlags().GetBool("version"); v {
		fmt.Printf("Version: %s, Built@%s\n", utils.Version, utils.BuildTime)
		os.Exit(0)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ryai")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
