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

type Grader struct {
	ID   int
	Name string
}

var KnownGraders = map[string]Grader{
	"C":       {ID: 1, Name: "C"},
	"OS":      {ID: 2, Name: "OS"},
	"C++":     {ID: 3, Name: "C++"},
	"Pascal":  {ID: 5, Name: "Pascal"},
	"Haskell": {ID: 6, Name: "Haskell"},
	"Python":  {ID: 7, Name: "Python"},
	"Java":    {ID: 8, Name: "Java"},
}
