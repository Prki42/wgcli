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

package submitCmd

import (
	"fmt"

	"github.com/Prki42/wgcli/config"
	"github.com/Prki42/wgcli/webgrade"
	"github.com/spf13/viper"
)

type Settings struct {
	auth    config.Auth
	problem config.Problem
}

var settings Settings

func guard(args []string) error {
	if len(args) == 0 {
		if !viper.IsSet("problem.sourceFile") {
			return fmt.Errorf("source file not specified")
		} else {
			settings.problem.SourceFile = viper.GetString("problem.sourceFile")
		}
	} else {
		settings.problem.SourceFile = args[0]
	}

	if viper.IsSet("problem.fileName") {
		settings.problem.FileName = viper.GetString("problem.fileName")
	} else {
		settings.problem.FileName = settings.problem.SourceFile
	}

	grader := viper.GetString("problem.grader")
	if _, exists := webgrade.KnownGraders[grader]; !exists {
		return fmt.Errorf("%s is not a known grader", grader)
	}
	settings.problem.Grader = grader

	id := viper.GetInt("problem.problemId")
	if id < 0 {
		return fmt.Errorf("problem id %d is not a valid id", id)
	}
	settings.problem.ProblemId = id

	return nil
}
