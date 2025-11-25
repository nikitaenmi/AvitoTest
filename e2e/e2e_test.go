package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) SetupSuite() {
	waitForService(s.T())
}

func (s *E2ETestSuite) Test01_CreateTeamAndPR() {
	t := s.T()

	teamName := generateUniqueID("team-1")
	user1 := generateUniqueID("user-1")
	user2 := generateUniqueID("user-2")
	user3 := generateUniqueID("user-3")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user1, Username: "Alice", IsActive: true},
			{UserID: user2, Username: "Bob", IsActive: true},
			{UserID: user3, Username: "Charlie", IsActive: true},
		},
	}
	createTeam(t, teamReq)

	prReq := CreatePRRequest{
		PullRequestID:   generateUniqueID("pr-1"),
		PullRequestName: "Test PR 1",
		AuthorID:        user1,
	}
	prResponse := createPR(t, prReq)

	assert.Equal(t, prReq.PullRequestID, prResponse.PR.PullRequestID)
	assert.Equal(t, "Test PR 1", prResponse.PR.PullRequestName)
	assert.Equal(t, user1, prResponse.PR.AuthorID)
	assert.Equal(t, "OPEN", prResponse.PR.Status)

	assert.NotEmpty(t, prResponse.PR.AssignedReviewers)
	for _, reviewer := range prResponse.PR.AssignedReviewers {
		assert.NotEqual(t, user1, reviewer)
	}
}

func (s *E2ETestSuite) Test02_ReassignReviewer() {
	t := s.T()

	teamName := generateUniqueID("team-2")
	user4 := generateUniqueID("user-4")
	user5 := generateUniqueID("user-5")
	user6 := generateUniqueID("user-6")
	user7 := generateUniqueID("user-7")
	user8 := generateUniqueID("user-8")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user4, Username: "David", IsActive: true},
			{UserID: user5, Username: "Eve", IsActive: true},
			{UserID: user6, Username: "Frank", IsActive: true},
			{UserID: user7, Username: "Grace", IsActive: true},
			{UserID: user8, Username: "Henry", IsActive: true},
		},
	}
	createTeam(t, teamReq)

	prID := generateUniqueID("pr-2")
	prReq := CreatePRRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR 2",
		AuthorID:        user4,
	}
	prResponse := createPR(t, prReq)

	require.NotEmpty(t, prResponse.PR.AssignedReviewers, "PR should have reviewers")
	assert.Len(t, prResponse.PR.AssignedReviewers, 2, "Should have exactly 2 reviewers")

	oldReviewer := prResponse.PR.AssignedReviewers[0]

	reassignReq := ReassignReviewerRequest{
		PullRequestID: prID,
		OldUserID:     oldReviewer,
	}
	reassignResponse := reassignReviewer(t, reassignReq)

	assert.NotEqual(t, oldReviewer, reassignResponse.ReplacedBy)
	assert.Contains(t, reassignResponse.PR.AssignedReviewers, reassignResponse.ReplacedBy)
	assert.NotContains(t, reassignResponse.PR.AssignedReviewers, oldReviewer)

	reviewers := reassignResponse.PR.AssignedReviewers
	reviewerSet := make(map[string]bool)
	for _, reviewer := range reviewers {
		assert.False(t, reviewerSet[reviewer], "Duplicate reviewer found: %s", reviewer)
		reviewerSet[reviewer] = true
	}
}

func (s *E2ETestSuite) Test03_ReassignPreventsDuplicates() {
	t := s.T()

	teamName := generateUniqueID("team-3")
	user9 := generateUniqueID("user-9")
	user10 := generateUniqueID("user-10")
	user11 := generateUniqueID("user-11")
	user12 := generateUniqueID("user-12")
	user13 := generateUniqueID("user-13")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user9, Username: "Ivan", IsActive: true},
			{UserID: user10, Username: "John", IsActive: true},
			{UserID: user11, Username: "Kate", IsActive: true},
			{UserID: user12, Username: "Leo", IsActive: true},
			{UserID: user13, Username: "Mona", IsActive: true},
		},
	}
	createTeam(t, teamReq)

	prID := generateUniqueID("pr-3")
	prReq := CreatePRRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR 3",
		AuthorID:        user9,
	}
	prResponse := createPR(t, prReq)

	require.Len(t, prResponse.PR.AssignedReviewers, 2, "PR should have 2 reviewers")

	oldReviewer := prResponse.PR.AssignedReviewers[0]
	secondReviewer := prResponse.PR.AssignedReviewers[1]

	for i := 0; i < 3; i++ {
		reassignReq := ReassignReviewerRequest{
			PullRequestID: prID,
			OldUserID:     oldReviewer,
		}
		reassignResponse := reassignReviewer(t, reassignReq)

		assert.NotEqual(t, secondReviewer, reassignResponse.ReplacedBy)

		assert.Contains(t, reassignResponse.PR.AssignedReviewers, secondReviewer)

		reviewers := reassignResponse.PR.AssignedReviewers
		assert.Len(t, reviewers, 2)

		oldReviewer = reassignResponse.ReplacedBy
	}
}

func (s *E2ETestSuite) Test04_CannotReassignOnMergedPR() {
	t := s.T()

	teamName := generateUniqueID("team-4")
	user14 := generateUniqueID("user-14")
	user15 := generateUniqueID("user-15")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user14, Username: "Oliver", IsActive: true},
			{UserID: user15, Username: "Penny", IsActive: true},
		},
	}
	createTeam(t, teamReq)

	prID := generateUniqueID("pr-merge")
	prReq := CreatePRRequest{
		PullRequestID:   prID,
		PullRequestName: "PR to Merge",
		AuthorID:        user14,
	}
	prResponse := createPR(t, prReq)

	require.NotEmpty(t, prResponse.PR.AssignedReviewers, "PR should have reviewers")
	oldReviewer := prResponse.PR.AssignedReviewers[0]

	mergePR(t, prID)

	reassignReq := ReassignReviewerRequest{
		PullRequestID: prID,
		OldUserID:     oldReviewer,
	}

	body, err := json.Marshal(reassignReq)
	require.NoError(t, err)

	baseURL := getBaseURL()
	resp, err := http.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, "PR_MERGED", errorResp.Error.Code)
}

func (s *E2ETestSuite) Test05_UserActivityAffectsAssignment() {
	t := s.T()

	teamName := generateUniqueID("team-5")
	user16 := generateUniqueID("user-16")
	user17 := generateUniqueID("user-17")
	user18 := generateUniqueID("user-18")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user16, Username: "Quinn", IsActive: true},
			{UserID: user17, Username: "Rachel", IsActive: true},
			{UserID: user18, Username: "Sam", IsActive: true},
		},
	}
	createTeam(t, teamReq)

	setUserActive(t, user17, false)

	prReq := CreatePRRequest{
		PullRequestID:   generateUniqueID("pr-activity"),
		PullRequestName: "PR Activity Test",
		AuthorID:        user16,
	}
	prResponse := createPR(t, prReq)

	for _, reviewer := range prResponse.PR.AssignedReviewers {
		assert.NotEqual(t, user17, reviewer, "Inactive user should not be assigned as reviewer")
	}
}

func (s *E2ETestSuite) Test06_ReassignWithLimitedAvailableReviewers() {
	t := s.T()

	teamName := generateUniqueID("team-small")
	user19 := generateUniqueID("user-19")
	user20 := generateUniqueID("user-20")

	smallTeamReq := TeamRequest{
		TeamName: teamName,
		Members: []UserRequest{
			{UserID: user19, Username: "Tina", IsActive: true},
			{UserID: user20, Username: "Ursula", IsActive: true},
		},
	}
	createTeam(t, smallTeamReq)

	prID := generateUniqueID("pr-small")
	prReq := CreatePRRequest{
		PullRequestID:   prID,
		PullRequestName: "Small Team PR",
		AuthorID:        user19,
	}
	prResponse := createPR(t, prReq)

	require.Len(t, prResponse.PR.AssignedReviewers, 1)
	oldReviewer := prResponse.PR.AssignedReviewers[0]

	reassignReq := ReassignReviewerRequest{
		PullRequestID: prID,
		OldUserID:     oldReviewer,
	}

	body, err := json.Marshal(reassignReq)
	require.NoError(t, err)

	baseURL := getBaseURL()
	resp, err := http.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, "NO_CANDIDATE", errorResp.Error.Code)
}
