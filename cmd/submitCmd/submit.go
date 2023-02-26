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
	"os"

	"github.com/Prki42/wgcli/cmd/utils"
	"github.com/Prki42/wgcli/webgrade"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger zerolog.Logger

func NewCommand() *cobra.Command {
	submitCmd := &cobra.Command{
		Use:               "submit",
		Short:             "Submit code to WebGrade",
		Long:              ``,
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: persistentPreRunE,
		PreRunE:           preRunE,
		Run:               run,
	}

	submitCmd.Flags().StringP("fileName", "f", "", "filename to be sent")
	viper.BindPFlag("problem.fileName", submitCmd.Flags().Lookup("problem.fileName"))

	submitCmd.Flags().IntP("problemId", "p", -1, "problem ID")
	viper.BindPFlag("problem.problemId", submitCmd.Flags().Lookup("problemId"))

	submitCmd.Flags().StringP("grader", "g", "C", "grader (language)")
	viper.BindPFlag("problem.grader", submitCmd.Flags().Lookup("grader"))

	return submitCmd
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	if err := utils.PersistPreRuns(cmd, args); err != nil {
		return err
	}

	logger = log.With().Str("command", "submit").Logger()

	return utils.RequireCredentials(&settings.auth)
}

func preRunE(cmd *cobra.Command, args []string) error {
	return guard(args)
}

func run(cmd *cobra.Command, args []string) {
	logger.Info().Msg("submit command called")
	logger.Debug().Interface("settings", settings.problem).Msg("submit parameters")

	// Submit
	client, err := webgrade.NewWebGradeClientLogin(settings.auth.Username, settings.auth.Password)
	if err != nil {
		logger.Error().Err(err).Msg("login")
	}
	cobra.CheckErr(err)
	defer client.Logout()

	f, err := os.Open(settings.problem.SourceFile)
	if err != nil {
		logger.Error().Err(err).Msg("opening source file")
	}
	cobra.CheckErr(err)
	defer f.Close()

	grader := settings.problem.Grader
	body := webgrade.SubmissionRequest{
		GraderID:   webgrade.KnownGraders[grader].ID,
		GraderName: grader,
		ProblemID:  settings.problem.ProblemId,
		FileName:   settings.problem.FileName,
	}
	logger.Debug().Interface("content", body).Msg("submission content")
	resp, err := client.SubmitCode(body, f)
	if err != nil {
		logger.Error().Err(err).Msg("submission response")
	}
	cobra.CheckErr(err)
	logger.Debug().Interface("response", resp).Msg("submission response")

	// Wait for task to finish and print test results
	prev := 0
	failedAttempts := 0
	finished := false
	for {
		task, err := client.GetTaskDetails(resp)
		if err != nil {
			logger.Warn().Err(err).Msg("fetching task details")
			failedAttempts++
			if failedAttempts >= 10 {
				fmt.Printf("Failed 10 consecutive times to check for task status\n")
				break
			}
			continue
		}
		logger.Debug().Interface("task", task).Msg("")
		length := len(task.Tests)
		if prev != length {
			for i := prev; i != length; i++ {
				if task.Tests[i].Output == 1 {
					fmt.Printf("(+) Test %d passed\n", i+1)
				} else {
					fmt.Printf("(-) Test %d failed\n", i+1)
				}
			}
			prev = length
		}
		if task.State == "finished" {
			finished = true
			break
		}
	}
	if !finished {
		fmt.Println("Code checking not complete")
	}

	// Retrieve submission details (score that got saved on webgrade)
	details, err := client.GetSubmissionDetails(resp.SubmissionID)
	if err != nil {
		logger.Error().Err(err).Msg("submission details")
	}
	cobra.CheckErr(err)
	logger.Debug().Interface("response", details).Msg("submission details")

	fmt.Printf("Score: %v\n", details.Score)
}
