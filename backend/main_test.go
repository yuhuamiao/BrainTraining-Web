package main

import (
	"braintraining/backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestApp(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	return setupRouter(db), db
}

func requestJSON(t *testing.T, app http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, req)
	return response
}

func registerAndLogin(t *testing.T, app http.Handler, username string) (string, string) {
	t.Helper()
	register := requestJSON(t, app, http.MethodPost, "/api/v1/register", map[string]string{
		"username": username,
		"password": "secret123",
	}, "")
	if register.Code != http.StatusOK {
		t.Fatalf("register returned %d: %s", register.Code, register.Body.String())
	}

	login := requestJSON(t, app, http.MethodPost, "/api/v1/login", map[string]string{
		"username": username,
		"password": "secret123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login returned %d: %s", login.Code, login.Body.String())
	}
	var response struct {
		Token string `json:"token"`
		User  struct {
			UserID string `json:"userId"`
		} `json:"user"`
	}
	if err := json.Unmarshal(login.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	return response.Token, response.User.UserID
}

func TestAuthAndScoresUseAuthenticatedUser(t *testing.T) {
	app, db := newTestApp(t)
	token, userID := registerAndLogin(t, app, "test-user")

	unauthorized := requestJSON(t, app, http.MethodPost, "/api/v1/bus/scores", map[string]any{
		"successNum": 3, "trainingNum": 4, "accuracy": 0.75, "level": "easy",
	}, "")
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("score without token returned %d", unauthorized.Code)
	}

	score := requestJSON(t, app, http.MethodPost, "/api/v1/bus/scores", map[string]any{
		"userId": "spoofed-user", "successNum": 3, "trainingNum": 4, "accuracy": 0.75, "level": "easy",
	}, token)
	if score.Code != http.StatusOK {
		t.Fatalf("submit score returned %d: %s", score.Code, score.Body.String())
	}

	var record models.BusRecord
	if err := db.First(&record).Error; err != nil {
		t.Fatalf("read score: %v", err)
	}
	if record.UserID != userID {
		t.Fatalf("score stored for %q, want authenticated user %q", record.UserID, userID)
	}

	list := requestJSON(t, app, http.MethodGet, "/api/v1/user/scores?userId=spoofed-user", nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("list scores returned %d: %s", list.Code, list.Body.String())
	}
	var scores map[string][]map[string]any
	if err := json.Unmarshal(list.Body.Bytes(), &scores); err != nil {
		t.Fatalf("decode scores: %v", err)
	}
	if len(scores["bus"]) != 1 {
		t.Fatalf("got %d bus records, want 1", len(scores["bus"]))
	}
}

func TestTrainingDataValidation(t *testing.T) {
	app, _ := newTestApp(t)

	memory := requestJSON(t, app, http.MethodGet, "/api/v1/memory/matrix?difficulty=easy", nil, "")
	if memory.Code != http.StatusOK {
		t.Fatalf("memory matrix returned %d: %s", memory.Code, memory.Body.String())
	}
	var matrix struct {
		Positions  []int `json:"positions"`
		TotalCells int   `json:"totalCells"`
	}
	if err := json.Unmarshal(memory.Body.Bytes(), &matrix); err != nil {
		t.Fatalf("decode memory matrix: %v", err)
	}
	if len(matrix.Positions) != 5 || matrix.TotalCells != 30 {
		t.Fatalf("memory matrix has %d targets/%d cells, want 5/30", len(matrix.Positions), matrix.TotalCells)
	}

	invalid := requestJSON(t, app, http.MethodGet, "/api/v1/schulte/matrix?size=11", nil, "")
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid Schulte size returned %d", invalid.Code)
	}
}
