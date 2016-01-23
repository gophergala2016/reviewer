// Copyright © 2016 See CONTRIBUTORS <alvaro.cabanas@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reviewer

import (
	reviewer "."
	"testing"
)

func mockConfigFileUsed() string {
	return ""
}

func mockConfigFileUsed2() string {
	return "myfile"
}

func mockIsSet(a string) bool {
	return false
}

func mockIsSet2(a string) bool {
	return true
}


type mockViper struct {}
func (v *mockViper) AllKeys()  {
	var a []string
	a[0] = "test"
	return a
}

func mockSub(a string) bool {
	return mockViper
}

type mockViper2 struct {}
func (v *mockViper2) AllKeys()  {
	var a []string
	a[0] = "test2"
	return a
}

func mockSub2(a string) bool {
	return mockViper
}

func mockGetString(a string) bool {
	if a =="repositories.test.username" {
		return "christofdamian"
	}else if a== "repositories.test.status" {
		return "true"
	}else if a=="repositories.test.required"{
		return "3"
	}else if a =="repositories.test2.username"{
		return ""
	}else if a== "repositories.test2.status"{
		return ""
	}else if a=="repositories.test2.required"{
		return ""
	}else{
		return ""
	}
}

func TestCheckFile(t *testing.T) {
	reviewer.ConfigFileUsed = mockConfigFileUsed

	err := reviewer.CheckFile()
	if err == nil {
		t.Fatal("Without configuration file or empty, it should complain")
	}

	reviewer.ConfigFileUsed = mockConfigFileUsed2

	err = reviewer.CheckFile()
	if err != nil {
		t.Fatal("With configuration file not empty, it should not complain")
	}
}

func TestCheckRepositories(t *testing.T) {
	reviewer.IsSet = mockIsSet

	err := reviewer.CheckRepositories()
	if err == nil {
		t.Fatal("Without Repositories tag or without any repository set, it should complain")
	}

	reviewer.IsSet = mockIsSet2

	err = reviewer.CheckRepositories()
	if err != nil {
		t.Fatal("With Repositories tag and repositories set, it should not complain")
	}
}

func TestCheckRepositoriesData(t *testing.T) {
	/*reviewer.Sub = mockSub
	reviewer.GetString = mockGetString

	response,err := reviewer.CheckRepositoriesData()
	if err != nil {
		t.Fatal("With Repositories tag and repositories set, it should not complain")
	}

	if response != {

	}

	reviewer.Sub = mockSub2

	response2,err2 := reviewer.CheckRepositoriesData()
	if err != nil {
		t.Fatal("With Repositories tag and repositories set, it should not complain")
	}*/
}