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

package utils

import (
	"fmt"
	"syscall"

	"github.com/Prki42/wgcli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func RequireCredentials(settings *config.Auth) error {
	if !viper.IsSet("auth.username") {
		return fmt.Errorf("username not provided")
	}
	settings.Username = viper.GetString("auth.username")

	if !viper.IsSet("auth.password") {
		fmt.Printf("Password for %v: ", settings.Username)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		cobra.CheckErr(err)
		fmt.Println()
		viper.Set("auth.password", string(bytePassword))
	}
	settings.Password = viper.GetString("auth.password")

	return nil
}
