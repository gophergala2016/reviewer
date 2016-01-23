// Copyright © 2016 See CONTRIBUTORS <ignasi.fosch@gmail.com>
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

	"net/url"
	"reflect"
	"testing"
)

// token contains the GH token.
var token = "GITHUB_USERS_TOKEN"

type mockGHClient struct {
	BaseURL *url.URL
}

func mockLookupEnv(k string) (string, bool) {
	if k == "REVIEWER_TOKEN" {
		return token, false
	}
	return "", true
}

func TestGetGHAuth(t *testing.T) {
	reviewer.LookupEnv = mockLookupEnv

	var result interface{}
	var errClient error
	result, errClient = reviewer.GetClient()

	if errClient != nil {
		t.Fatalf("GetClient returned error(%s) when everything was ok", errClient)
	}
	v, err := result.(mockGHClient)
	if err {
		t.Fatalf("GetClient returned %s instead of github.Client", reflect.TypeOf(v))
	}
}

func TestCommentSuccessScore(t *testing.T) {

	testScore := func(comment string, expected int) {
		score := getCommentSuccessScore(comment)
		if expected != score {
			t.Fatalf("Bad score %v (expected %v) for comment %v", score, expected, comment)
		}
	}

	testScore("Don't do it", 0)
	testScore("Yes +1", 1)
	testScore(":+1", 1)
	testScore("-1", -1)
	testScore("Oops +1 :-1: +1", 0)
}
