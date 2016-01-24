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

type mockConfig struct {
	key string
}

func (v *mockConfig) AllKeys() []string {
	a := []string{v.key}
	return a
}

func (v *mockConfig) GetString(key string) string {
	if key == "test.username" {
		return "christofdamian"
	} else if key == "test.status" {
		return "true"
	} else if key == "test.required" {
		return "3"
	} else if key == "test2.username" {
		return ""
	} else if key == "test2.status" {
		return ""
	} else if key == "test2.required" {
		return ""
	} else {
		return ""
	}
}

func newMockConfig(key string) *mockConfig {
	return &mockConfig{
		key: key,
	}
}

func mockConfigFileUsed() string {
	return ""
}

func mockConfigFileUsed2() string {
	return "/my/path/to/config/file"
}

func mockIsSet(a string) bool {
	return false
}

func mockIsSet2(a string) bool {
	return true
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
	config := newMockConfig("test")

	response, err := reviewer.CheckRepositoriesData(config)
	if err != nil {
		t.Fatal("With parameters set not return error")
	}

	if response == "- christofdamian / test ENABLED +1:3" {
		t.Fatal("With Repositories tag and repositories set, it should not complain")
	}

	config = newMockConfig("test2")

	response, err = reviewer.CheckRepositoriesData(config)
	if err == nil {
		t.Fatal("With parameters not set it should return error")
	}

	if response != "" {
		t.Fatal("With Repositories tag and repositories not set, it should complain")
	}

}
