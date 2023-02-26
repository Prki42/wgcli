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
	"path/filepath"

	"github.com/Prki42/wgcli/cmd/submitCmd"
	"github.com/Prki42/wgcli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func NewCommand() *cobra.Command {
	wgcli := &cobra.Command{
		Use:              "wgcli",
		Short:            "WebGrade CLI utility",
		Long:             ``,
		PersistentPreRun: persistentPreRun,
	}

	config.LoadGlobalConfig()

	confPath, _ := filepath.Abs("./")
	config.LoadConfig(confPath, ".wgcli.yaml")

	createCommandTree(wgcli)

	return wgcli
}

func createCommandTree(cmd *cobra.Command) {
	cmd.AddCommand(submitCmd.NewCommand())

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "additional config file")
	cmd.PersistentFlags().StringP("username", "U", "", "username")
	cmd.PersistentFlags().StringP("password", "P", "", "password")
	viper.BindPFlag("auth.username", cmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("auth.password", cmd.PersistentFlags().Lookup("password"))
}

func persistentPreRun(cnd *cobra.Command, args []string) {
	if cfgFile == "" {
		return
	}
	confPath, _ := filepath.Abs(cfgFile)
	fileName := filepath.Base(confPath)
	dirPath := filepath.Dir(confPath)
	config.LoadConfig(dirPath, fileName)
}
