/*
   Copyright 2023 Aleksa Prtenjača <aleksa.prtenjaca03@gmail.com>

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
	"path/filepath"

	"github.com/Prki42/wgcli/cmd/submitCmd"
	"github.com/Prki42/wgcli/config"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var verbose bool
var logFile string

func NewCommand() *cobra.Command {
	wgcli := &cobra.Command{
		Use:              "wgcli",
		Short:            "WebGrade CLI utility",
		Long:             ``,
		PersistentPreRun: persistentPreRun,
	}

	createCommandTree(wgcli)

	return wgcli
}

func createCommandTree(cmd *cobra.Command) {
	cmd.AddCommand(submitCmd.NewCommand())

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "additional config file")
	cmd.PersistentFlags().StringP("username", "U", "", "username")
	cmd.PersistentFlags().StringP("password", "P", "", "password")
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "prints debug info")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "logs to stderr")
	cmd.PersistentFlags().StringVarP(&logFile, "logFile", "l", "", "different log file")
	cmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	viper.BindPFlag("auth.username", cmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("auth.password", cmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("noColor", cmd.PersistentFlags().Lookup("no-color"))
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	setupLogging()

	// Global config
	globalConfPath, err := config.LoadGlobalConfig()
	if err != nil {
		log.Error().Err(err).Str("file", globalConfPath).Msg("loading global log file")
	} else {
		log.Info().Str("file", globalConfPath).Msg("global config file loaded")
	}

	// Local config (./)
	localConfPath, _ := filepath.Abs("./")
	localConfPath, err = config.LoadConfig(localConfPath, ".wgcli.yaml")
	if err != nil {
		log.Error().Err(err).Str("file", localConfPath).Msg("loading local log file")
	} else {
		log.Info().Str("file", localConfPath).Msg("local config file loaded")
	}

	// Config file passed as argument
	if cfgFile != "" {
		confPath, err := filepath.Abs(cfgFile)
		if err != nil {
			log.Error().Err(err).Str("file", cfgFile).Msg("loading custom log file")
			return
		}
		fileName := filepath.Base(confPath)
		dirPath := filepath.Dir(confPath)
		confPath, err = config.LoadConfig(dirPath, fileName)
		if err != nil {
			log.Error().Err(err).Str("file", confPath).Msg("laoding custom log file")
			return
		}
		log.Info().Str("file", cfgFile).Msg("custom log file loaded")
	}

	color.NoColor = viper.GetBool("noColor")
}

func setupLogging() {
	if logFile == "" {
		confPath, err := os.UserConfigDir()
		if err != nil && verbose {
			fmt.Println("Error opening log file:", err)
		}
		logFile = filepath.Join(confPath, "wgcli", "log.json")
	} else {
		var err error
		logFile, err = filepath.Abs(logFile)
		if err != nil && verbose {
			fmt.Println("Error getting log file path:", err)
		}
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil && verbose {
		fmt.Println("Error opening log file:", err)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if verbose {
		multi := zerolog.MultiLevelWriter(file, zerolog.ConsoleWriter{Out: os.Stderr})
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(file).With().Timestamp().Logger()
	}
}
