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
	"errors"
	"strings"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// GetString contains the function used to lookup environment variables.
var GetString = viper.GetString

// ChangesServicer is an interface for listing changes.
type ChangesServicer interface {
	List(string, string, *github.PullRequestListOptions) ([]github.PullRequest, *github.Response, error)
}

// TicketsServicer is an interface for listing changes.
type TicketsServicer interface {
	ListComments(string, string, int, *github.IssueListCommentsOptions) ([]github.IssueComment, *github.Response, error)
}

// GHClient is the wrapper around github.Client.
type GHClient struct {
	client *github.Client
	Changes ChangesServicer
	Tickets TicketsServicer
}

// NewGHClient is the constructor for GHClient.
func NewGHClient(httpClient *http.Client) *GHClient {
	client := &GHClient{
		client: github.NewClient(httpClient),
	}
	client.Changes = client.client.PullRequests
	client.Tickets = client.client.Issues
	return client
}

// PullRequestInfo contains the id, title, and CR score of a pull request.
type PullRequestInfo struct {
	Number int // id of the pull request
	Title  string
	Score  int
}

// GetClient returns a github.Client authenticated.
func GetClient() (*GHClient, error) {
	token := GetString("authorization.token")
	if token == "" {
		return nil, errors.New("An error occurred getting REVIEWER_TOKEN environment variable\n")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return NewGHClient(tc), nil
}

// getCommentSuccesScore returns the score for the Comment.
func getCommentSuccessScore(comment string) int {
	score := 0
	if strings.Contains(comment, "+1") {
		score++
	}
	if strings.Contains(comment, "-1") {
		score--
	}
	return score
}

// GetPullRequestInfos returns the list of pull requests and the CR success score based on comments
func GetPullRequestInfos(client *GHClient, owner string, repo string) ([]PullRequestInfo, error) {
	//TODO: At this moment if there's a lot of PR, does not returns the full list, needs pagination.
	//      Also maybe we need to take care about how much requests are done in order to not exceed
	//      the quota.

	pullRequests, _, err := client.Changes.List(owner, repo, nil)
	if err != nil {
		return nil, err
	}
	pris := make([]PullRequestInfo, len(pullRequests))
	for n, pullRequest := range pullRequests {
		pris[n].Number = *pullRequest.Number
		pris[n].Title = *pullRequest.Title
		comments, _, err := client.Tickets.ListComments(owner, repo, *pullRequest.Number, nil)
		if err != nil {
			return nil, err
		}
		for _, comment := range comments {
			if comment.Body == nil {
				continue
			}
			pris[n].Score += getCommentSuccessScore(*comment.Body)
		}
	}
	return pris, nil
}

// IsMergeable returns true if the PullRequest is mergeable.
func IsMergeable(pullRequest *github.PullRequest) bool {
	return (pullRequest.Mergeable != nil) && (*pullRequest.Mergeable)
}

// PassedTests checks if the PR statuses are ok.
func PassedTests(client *GHClient, pullRequest *github.PullRequest, owner string, repo string) (bool, error) {
	head := *pullRequest.Head.SHA
	combinedStatus, _, err := client.client.Repositories.GetCombinedStatus(owner, repo, head, nil)

	if err != nil {
		return false, err
	}
	return (*combinedStatus.State != "success"), nil
}
