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
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type WebGradeClient struct {
	client http.Client
}

func NewWebGradeClientLogin(username, password string) (*WebGradeClient, error) {
	wgClient := &WebGradeClient{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	wgClient.client = http.Client{
		Jar: jar,
	}

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	params.Add("action", "login")

	resp, err := wgClient.client.PostForm(baseUrl+authEndpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if string(respData) == "null" {
		return nil, fmt.Errorf("login failed, bad credentials")
	}

	return wgClient, nil
}

func (w WebGradeClient) Logout() error {
	_, err := w.client.Get(baseUrl + authEndpoint + "?action=logout")
	if err != nil {
		return err
	}
	return nil
}
