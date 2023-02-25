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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	return utils.RequireCredentials(&settings.auth)
}

func preRunE(cmd *cobra.Command, args []string) error {
	return guard(args)
}

func run(cmd *cobra.Command, args []string) {
	client, err := webgrade.NewWebGradeClientLogin(settings.auth.Username, settings.auth.Password)
	cobra.CheckErr(err)
	defer client.Logout()

	f, err := os.Open(settings.problem.SourceFile)
	cobra.CheckErr(err)
	defer f.Close()

	grader := settings.problem.Grader
	resp, err := client.SubmitCode(webgrade.SubmissionRequest{
		GraderID:   webgrade.KnownGraders[grader].ID,
		GraderName: grader,
		ProblemID:  settings.problem.ProblemId,
		FileName:   settings.problem.FileName,
	}, f)
	cobra.CheckErr(err)

	prev := -1
	for {
		task, err := client.GetTaskDetails(resp)
		length := len(task.Tests)
		if prev != length {
			prev = length
			if length > 0 {
				fmt.Printf("Finished %v test(s)\n", length)
			}
		}
		if err != nil || task.State == "finished" {
			break
		}
	}
	if err != nil {
		fmt.Printf("Failed checking if all tasks are done\n")
	}
	err = nil

	details, err := client.GetSubmissionDetails(resp.SubmissionID)
	cobra.CheckErr(err)

	fmt.Printf("Score: %v\n", details.Score)
}
