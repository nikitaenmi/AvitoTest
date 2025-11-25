package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getBaseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

type UserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamRequest struct {
	TeamName string        `json:"team_name"`
	Members  []UserRequest `json:"members"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type PRResponse struct {
	PR struct {
		PullRequestID     string   `json:"pull_request_id"`
		PullRequestName   string   `json:"pull_request_name"`
		AuthorID          string   `json:"author_id"`
		Status            string   `json:"status"`
		AssignedReviewers []string `json:"assigned_reviewers"`
	} `json:"pr"`
	ReplacedBy string `json:"replaced_by,omitempty"`
}

type TeamResponse struct {
	Team struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	} `json:"team"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func waitForService(t *testing.T) {
	baseURL := getBaseURL()
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("Service is not available")
}

func createTeam(t *testing.T, teamReq TeamRequest) {
	baseURL := getBaseURL()
	body, err := json.Marshal(teamReq)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		assert.Fail(t, "Failed to create team", "Expected 201 or 400, got %d", resp.StatusCode)
	}
}

func createPR(t *testing.T, prReq CreatePRRequest) *PRResponse {
	baseURL := getBaseURL()
	body, err := json.Marshal(prReq)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		assert.Fail(t, "Failed to create PR", "Expected 201 or 409, got %d", resp.StatusCode)
	}

	var prResponse PRResponse
	err = json.NewDecoder(resp.Body).Decode(&prResponse)
	if err != nil && resp.StatusCode == http.StatusConflict {
		return &prResponse
	}
	require.NoError(t, err)

	return &prResponse
}

func reassignReviewer(t *testing.T, req ReassignReviewerRequest) *PRResponse {
	baseURL := getBaseURL()
	body, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err == nil {
			t.Logf("Reassign failed: %s - %s", errorResp.Error.Code, errorResp.Error.Message)
		}
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to reassign reviewer")

	var prResponse PRResponse
	err = json.NewDecoder(resp.Body).Decode(&prResponse)
	require.NoError(t, err)

	return &prResponse
}

func setUserActive(t *testing.T, userID string, isActive bool) {
	baseURL := getBaseURL()
	req := map[string]interface{}{
		"user_id":   userID,
		"is_active": isActive,
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/users/setIsActive", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to set user active status")
}

func getTeam(t *testing.T, teamName string) *TeamResponse {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/team/get?team_name=" + teamName)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to get team")

	var teamResponse TeamResponse
	err = json.NewDecoder(resp.Body).Decode(&teamResponse)
	require.NoError(t, err)

	return &teamResponse
}

func mergePR(t *testing.T, prID string) {
	baseURL := getBaseURL()
	req := map[string]interface{}{
		"pull_request_id": prID,
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/pullRequest/merge", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to merge PR")
}

func generateUniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
