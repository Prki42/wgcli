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

package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Reads default config file
func LoadGlobalConfig() error {
	confPath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	confPath = filepath.Join(confPath, "wgcli")

	viper.AddConfigPath(confPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("conf")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// Merges global config with additional config file (overwrites)
func LoadConfig(filepath, filename string) error {
	v := viper.New()

	v.AddConfigPath(filepath)
	v.SetConfigType("yaml")
	v.SetConfigName(filename)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
		return err
	}

	return nil
}
