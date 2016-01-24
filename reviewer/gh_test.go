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

	"github.com/google/go-github/github"
	"reflect"
	"testing"
)

// token contains the GH token.
var token = "GITHUB_USERS_TOKEN"

// mockChangesService is a mock for github.PullRequestsService.
type mockChangesService struct {
	listPullRequests []github.PullRequest
}

// newMockChangesService creates a new ChangesService implementation.
func newMockChangesService(listPR []github.PullRequest) *mockChangesService {
	return &mockChangesService{
		listPullRequests: listPR,
	}
}

// mockChangesService's List implementation.
func (m *mockChangesService) List(owner string, repo string, opt *github.PullRequestListOptions) ([]github.PullRequest, *github.Response, error) {
	return m.listPullRequests, nil, nil
}

// mockTicketsService is a mock for github.PullRequestsService.
type mockTicketsService struct {
	listIssueComments [][]github.IssueComment
}

// newMockTicketsService creates a new TicketsService implementation.
func newMockTicketsService(listIssueComments [][]github.IssueComment) *mockTicketsService {
	return &mockTicketsService{
		listIssueComments: listIssueComments,
	}
}

// mockTicketsService's List implementation.
func (m *mockTicketsService) ListComments(owner string, repo string, number int, opt *github.IssueListCommentsOptions) ([]github.IssueComment, *github.Response, error) {
	return nil, nil, nil
}

// Constructor for mockGHClient.
func newMockGHClient(listPR []github.PullRequest, listIssueComments [][]github.IssueComment) *reviewer.GHClient {
	client := &reviewer.GHClient{}
	client.Changes = newMockChangesService(listPR)
	client.Tickets = newMockTicketsService(listIssueComments)
	return client
}

func mockGetString(k string) string {
	if k == "authorization.token" {
		return token
	}
	return ""
}

func TestGetGHAuth(t *testing.T) {
	reviewer.GetString = mockGetString

	var result interface{}
	var errClient error
	result, errClient = reviewer.GetClient()

	if errClient != nil {
		t.Fatalf("GetClient returned error(%s) when everything was ok", errClient)
	}
	v, err := result.(reviewer.GHClient)
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

func newMockPullRequest(number int, title string) github.PullRequest {
	return github.PullRequest{
		Number: &number,
		Title:  &title,
	}
}

func TestGetPullRequestsInfo(t *testing.T) {
	//TODO: https://github.com/gophergala2016/reviewer/issues/22
	var emptyListPR []github.PullRequest
	emptyListPR = make([]github.PullRequest, 0)
	var emptyListIC [][]github.IssueComment
	emptyListIC = make([][]github.IssueComment, 0)
	client := newMockGHClient(emptyListPR, emptyListIC)

	var result []reviewer.PullRequestInfo
	var err error
	result, err = reviewer.GetPullRequestInfos(client, "user", "repo")

	if err != nil {
		t.Fatalf("Something went wrong when getting PR information")
	}
	if len(result) != 0 {
		t.Fatal("Got a populated list of PRInfos")
	}

	onePR := make([]github.PullRequest, 1)
	onePR[0] = newMockPullRequest(10, "Initial PR")
	client = newMockGHClient(onePR, emptyListIC)

	result, err = reviewer.GetPullRequestInfos(client, "user", "repo")

	if err != nil {
		t.Fatalf("Something went wrong when getting PR information")
	}
	if len(result) != 1 {
		t.Fatal("Got a incorrect quantity of PRInfos:", len(result))
	}

	twoPR := make([]github.PullRequest, 2)
	twoPR[0] = newMockPullRequest(10, "Initial PR")
	twoPR[1] = newMockPullRequest(11, "Not so initial PR")
	client = newMockGHClient(twoPR, emptyListIC)

	result, err = reviewer.GetPullRequestInfos(client, "user", "repo")

	if err != nil {
		t.Fatalf("Something went wrong when getting PR information")
	}
	if len(result) != 2 {
		t.Fatal("Got a incorrect quantity of PRInfos:", len(result))
	}
}
