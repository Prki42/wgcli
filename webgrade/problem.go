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

import "fmt"

type ProblemDetailsResult struct {
	BreadcrumbData struct {
		Course struct {
			Name string `json:"name"`
		} `json:"course"`
		Topic   string `json:"topic"`
		Problem struct {
			Name string `json:"name"`
		} `json:"problem"`
	} `json:"breadcrumbData"`
	AuthData struct {
		Username   string `json:"username"`
		UserID     int    `json:"userId"`
		Email      string `json:"email"`
		RoleName   string `json:"roleName"`
		IsVerified int    `json:"isVerified"`
	} `json:"authData"`
	ProblemSubmissions []struct {
		Score           int         `json:"score"`
		Timestamp       string      `json:"timestamp"`
		RequestedReview int         `json:"requestedReview"`
		RequestText     interface{} `json:"requestText"`
		SourceFile      string      `json:"sourceFile"`
		CodeReview      interface{} `json:"codeReview"`
		Reviewer        interface{} `json:"reviewer"`
		SubmissionID    int         `json:"submissionId"`
		ReviewerAvatar  interface{} `json:"reviewerAvatar"`
		ReviewerID      interface{} `json:"reviewerId"`
		SeenReview      interface{} `json:"seenReview"`
		Extension       string      `json:"extension"`
	} `json:"problemSubmissions"`
	ProblemDetails struct {
		ProblemDetails struct {
			Name                string      `json:"name"`
			Text                string      `json:"text"`
			TimeLimit           int         `json:"timeLimit"`
			MemoryLimit         int         `json:"memoryLimit"`
			TestCasesNum        int         `json:"testCasesNum"`
			TemplateFileContent interface{} `json:"templateFileContent"`
		} `json:"problemDetails"`
		Graders []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"graders"`
	} `json:"problemDetails"`
	Professors []struct {
		IDUser    int    `json:"idUser"`
		Avatar    string `json:"avatar"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"professors"`
	Reports         []interface{} `json:"reports"`
	CanSubmitReport bool          `json:"canSubmitReport"`
	CompanyData     interface{}   `json:"companyData"`
}

type ProblemData struct {
	CourseID  int
	TopicID   int
	ProblemID int
}

func (w WebGradeClient) GetProblemDetails(problem ProblemData) (ProblemDetailsResult, error) {
	topicId := "null"
	courseId := "null"
	if problem.CourseID > 0 {
		courseId = fmt.Sprint(problem.CourseID)
	}
	if problem.TopicID > 0 {
		topicId = fmt.Sprint(problem.TopicID)
	}
	urlEnd := fmt.Sprintf("?action=getProblemDetailsComponentData&courseId=%s&topicId=%s&problemId=%d", courseId, topicId, problem.ProblemID)
	resp, err := w.client.Get(baseUrl + studentEndpoint + urlEnd)
	if err != nil {
		return ProblemDetailsResult{}, err
	}
	defer resp.Body.Close()

	var result ProblemDetailsResult
	err = unmarshalBody(resp, &result)
	if err != nil {
		return ProblemDetailsResult{}, err
	}

	return result, nil
}
