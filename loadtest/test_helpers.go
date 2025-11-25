package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type PRResponse struct {
	PR struct {
		PullRequestID     string   `json:"pull_request_id"`
		PullRequestName   string   `json:"pull_request_name"`
		AuthorID          string   `json:"author_id"`
		Status            string   `json:"status"`
		AssignedReviewers []string `json:"assigned_reviewers"`
	} `json:"pr"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Request struct {
	Method   string
	URL      string
	Body     interface{}
	Expected int
}

var (
	userPool  []string
	teamPool  []string
	poolMutex sync.RWMutex
)

func generateUniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func initializeTestData(cfg *Config) {
	for i := 0; i < 5; i++ {
		teamName := generateUniqueID("load-team")
		user1 := generateUniqueID("user")
		user2 := generateUniqueID("user")
		user3 := generateUniqueID("user")

		teamReq := TeamRequest{
			TeamName: teamName,
			Members: []UserRequest{
				{UserID: user1, Username: fmt.Sprintf("User%d-1", i), IsActive: true},
				{UserID: user2, Username: fmt.Sprintf("User%d-2", i), IsActive: true},
				{UserID: user3, Username: fmt.Sprintf("User%d-3", i), IsActive: true},
			},
		}

		if err := sendRequest(cfg, Request{
			Method:   "POST",
			URL:      cfg.BaseURL + "/team/add",
			Body:     teamReq,
			Expected: http.StatusCreated,
		}, http.StatusCreated); err == nil {
			poolMutex.Lock()
			teamPool = append(teamPool, teamName)
			userPool = append(userPool, user1, user2, user3)
			poolMutex.Unlock()
		}
	}
}

func generateRequest(id int, baseURL string) (string, Request, int) {
	poolMutex.RLock()
	defer poolMutex.RUnlock()

	hasUsers := len(userPool) > 0
	hasTeams := len(teamPool) > 0

	switch id % 8 {
	case 0:
		teamName := generateUniqueID("team")
		user1 := generateUniqueID("user")
		user2 := generateUniqueID("user")

		teamReq := TeamRequest{
			TeamName: teamName,
			Members: []UserRequest{
				{UserID: user1, Username: "Member1", IsActive: true},
				{UserID: user2, Username: "Member2", IsActive: true},
			},
		}
		return "create_team", Request{
			Method:   "POST",
			URL:      baseURL + "/team/add",
			Body:     teamReq,
			Expected: http.StatusCreated,
		}, http.StatusCreated

	case 1:
		if hasUsers {
			authorID := userPool[id%len(userPool)]
			prReq := CreatePRRequest{
				PullRequestID:   generateUniqueID("pr"),
				PullRequestName: "Load Test PR",
				AuthorID:        authorID,
			}
			return "create_pr", Request{
				Method:   "POST",
				URL:      baseURL + "/pullRequest/create",
				Body:     prReq,
				Expected: http.StatusCreated,
			}, http.StatusCreated
		}
		return "health_check", Request{
			Method:   "GET",
			URL:      baseURL + "/health",
			Body:     nil,
			Expected: http.StatusOK,
		}, http.StatusOK

	case 2:
		if hasUsers {
			userID := userPool[id%len(userPool)]
			setActiveReq := SetUserActiveRequest{
				UserID:   userID,
				IsActive: id%2 == 0,
			}
			return "set_user_active", Request{
				Method:   "POST",
				URL:      baseURL + "/users/setIsActive",
				Body:     setActiveReq,
				Expected: http.StatusOK,
			}, http.StatusOK
		}
		return "health_check", Request{
			Method:   "GET",
			URL:      baseURL + "/health",
			Body:     nil,
			Expected: http.StatusOK,
		}, http.StatusOK

	case 3:
		return "health_check", Request{
			Method:   "GET",
			URL:      baseURL + "/health",
			Body:     nil,
			Expected: http.StatusOK,
		}, http.StatusOK

	case 4:
		if hasTeams {
			teamName := teamPool[id%len(teamPool)]
			return "get_team", Request{
				Method:   "GET",
				URL:      baseURL + "/team/get?team_name=" + teamName,
				Body:     nil,
				Expected: http.StatusOK,
			}, http.StatusOK
		}
		return "health_check", Request{
			Method:   "GET",
			URL:      baseURL + "/health",
			Body:     nil,
			Expected: http.StatusOK,
		}, http.StatusOK

	case 5:
		if hasTeams {
			teamName := teamPool[0]
			teamReq := TeamRequest{
				TeamName: teamName,
				Members: []UserRequest{
					{UserID: generateUniqueID("user"), Username: "DuplicateUser", IsActive: true},
				},
			}
			return "create_team_duplicate", Request{
				Method:   "POST",
				URL:      baseURL + "/team/add",
				Body:     teamReq,
				Expected: http.StatusBadRequest,
			}, http.StatusBadRequest
		}
		return "health_check", Request{
			Method:   "GET",
			URL:      baseURL + "/health",
			Body:     nil,
			Expected: http.StatusOK,
		}, http.StatusOK

	case 6:
		prReq := CreatePRRequest{
			PullRequestID:   generateUniqueID("pr"),
			PullRequestName: "Invalid PR",
			AuthorID:        "non-existent-user-12345",
		}
		return "create_pr_invalid", Request{
			Method:   "POST",
			URL:      baseURL + "/pullRequest/create",
			Body:     prReq,
			Expected: http.StatusNotFound,
		}, http.StatusNotFound

	case 7:
		setActiveReq := SetUserActiveRequest{
			UserID:   "non-existent-user-12345",
			IsActive: true,
		}
		return "set_user_active_invalid", Request{
			Method:   "POST",
			URL:      baseURL + "/users/setIsActive",
			Body:     setActiveReq,
			Expected: http.StatusNotFound,
		}, http.StatusNotFound
	}

	return "health_check", Request{
		Method:   "GET",
		URL:      baseURL + "/health",
		Body:     nil,
		Expected: http.StatusOK,
	}, http.StatusOK
}

func sendRequest(cfg *Config, request Request, expectedStatus int) error {
	var bodyBytes []byte
	var err error

	if request.Body != nil {
		bodyBytes, err = json.Marshal(request.Body)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
	}

	client := &http.Client{Timeout: cfg.Timeout}
	var httpReq *http.Request

	if request.Body != nil {
		httpReq, err = http.NewRequest(request.Method, request.URL, bytes.NewReader(bodyBytes))
	} else {
		httpReq, err = http.NewRequest(request.Method, request.URL, nil)
	}

	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != expectedStatus && resp.StatusCode < 500 {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("status %d: %s", resp.StatusCode, errorResp.Error.Message)
		}
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func isExpectedError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "status 404") ||
		strings.Contains(errStr, "status 400") ||
		strings.Contains(errStr, "status 409")
}
