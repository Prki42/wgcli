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

package webgrade

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type SubmissionRequest struct {
	GraderID   int
	ProblemID  int
	GraderName string
	FileName   string
}

type SubmissionResponse struct {
	GraderTaskID string `json:"graderTaskId"`
	SubmissionID string `json:"submissionId"`
}

type TaskStatusResponse struct {
	State string `json:"state"`
	Tests []struct {
		Status   string `json:"status"`
		Code     int    `json:"code"`
		Output   int    `json:"output"`
		Time     int    `json:"time"`
		WallTime int    `json:"wall_time"`
		Memory   int    `json:"memory"`
	} `json:"tests,omitempty"`
}

type SubmissionDetailsResponse struct {
	Score           string      `json:"score"`
	Timestamp       string      `json:"timestamp"`
	RequestedReview string      `json:"requestedReview"`
	RequestText     interface{} `json:"requestText"`
	SourceFile      string      `json:"sourceFile"`
	CodeReview      interface{} `json:"codeReview"`
	Reviewer        interface{} `json:"reviewer"`
	SubmissionID    string      `json:"submissionId"`
	ReviewerAvatar  interface{} `json:"reviewerAvatar"`
	ReviewerID      interface{} `json:"reviewerId"`
	SeenReview      string      `json:"seenReview"`
	Extension       string      `json:"extension"`
}

func (w WebGradeClient) SubmitCode(requestData SubmissionRequest, source io.Reader) (SubmissionResponse, error) {
	formBody := &bytes.Buffer{}
	formDataContentType := fillFormBody(formBody, requestData, source)

	req, err := http.NewRequest("POST", baseUrl+submissionEndpoint, formBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", formDataContentType)

	resp, err := w.client.Do(req)
	if err != nil {
		return SubmissionResponse{}, err
	}
	defer resp.Body.Close()

	var result SubmissionResponse
	err = unmarshalBody(resp, &result)

	if err != nil {
		return SubmissionResponse{}, err
	}

	return result, nil
}

func (w WebGradeClient) GetSubmissionDetails(submissionId string) (SubmissionDetailsResponse, error) {
	resp, err := w.client.Get(baseUrl + submissionEndpoint + "?action=getSubmissionDetails&submissionId=" + submissionId)
	if err != nil {
		return SubmissionDetailsResponse{}, err
	}
	defer resp.Body.Close()

	var result SubmissionDetailsResponse
	err = unmarshalBody(resp, &result)
	if err != nil {
		return SubmissionDetailsResponse{}, err
	}

	return result, nil
}

func (w WebGradeClient) GetTaskDetails(submission SubmissionResponse) (TaskStatusResponse, error) {
	resp, err := w.client.Get(baseUrl + submissionEndpoint + "?action=getTaskStatus&graderTaskId=" + submission.GraderTaskID + "&submissionId=" + submission.SubmissionID)
	if err != nil {
		return TaskStatusResponse{}, err
	}
	defer resp.Body.Close()

	var result TaskStatusResponse
	err = unmarshalBody(resp, &result)
	if err != nil {
		return TaskStatusResponse{}, err
	}

	return result, nil
}

func fillFormBody(formBody *bytes.Buffer, requestData SubmissionRequest, source io.Reader) string {
	fileWriter := multipart.NewWriter(formBody)

	fileWriter.WriteField("action", "postSubmission")
	fileWriter.WriteField("problemId", fmt.Sprint(requestData.ProblemID))
	fileWriter.WriteField("graderId", fmt.Sprint(requestData.GraderID))
	fileWriter.WriteField("graderName", fmt.Sprint(requestData.GraderName))
	part, _ := fileWriter.CreateFormFile("file", requestData.FileName)
	io.Copy(part, source)
	fileWriter.WriteField("template", "null")
	fileWriter.Close()

	return fileWriter.FormDataContentType()
}
