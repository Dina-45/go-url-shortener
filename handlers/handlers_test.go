package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-url-shortener/database"
	"go-url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	db.AutoMigrate(&models.User{}, &models.URL{})
	database.DB = db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestRegister_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/register", Register)

	body := `{"username":"alice","password":"secret123"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/register", Register)

	body := `{"username":"bob","password":"pass"}`
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if i == 1 {
			assert.Equal(t, http.StatusConflict, w.Code)
		}
	}
}
func TestLogin_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/register", Register)
	r.POST("/login", Login)

	reg := `{"username":"carol","password":"mypass"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(reg))
	req.Header.Set("Content-Type", "application/json")
	httptest.NewRecorder()
	r.ServeHTTP(httptest.NewRecorder(), req)

	login := `{"username":"carol","password":"mypass"}`
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(login))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req2)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/register", Register)
	r.POST("/login", Login)

	reg := `{"username":"dave","password":"correct"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(reg))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	login := `{"username":"dave","password":"wrong"}`
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(login))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req2)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateURL_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()

	r.POST("/urls", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		CreateURL(c)
	})

	body := `{"original_url":"https://google.com","short_code":"ggl"}`
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ggl", resp["short_code"])
}

func TestCreateURL_DuplicateCode(t *testing.T) {
	setupTestDB()
	r := setupRouter()

	r.POST("/urls", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		CreateURL(c)
	})

	body := `{"original_url":"https://example.com","short_code":"dup"}`
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("POST", "/urls", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if i == 1 {
			assert.Equal(t, http.StatusConflict, w.Code)
		}
	}
}

func TestGetAllURLs_Success(t *testing.T) {
	setupTestDB()

	database.DB.Create(&models.URL{
		OriginalURL: "https://ya.ru",
		ShortCode:   "ya",
		UserID:      42,
	})

	r := setupRouter()
	r.GET("/urls", func(c *gin.Context) {
		c.Set("user_id", uint(42))
		GetAllURLs(c)
	})

	req, _ := http.NewRequest("GET", "/urls", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var urls []models.URL
	json.Unmarshal(w.Body.Bytes(), &urls)
	assert.Len(t, urls, 1)
	assert.Equal(t, "ya", urls[0].ShortCode)
}

func TestRedirectShortURL_Success(t *testing.T) {
	setupTestDB()

	database.DB.Create(&models.URL{
		OriginalURL: "https://openai.com",
		ShortCode:   "oai",
		UserID:      1,
	})

	r := setupRouter()
	r.GET("/short/:code", RedirectShortURL)

	req, _ := http.NewRequest("GET", "/short/oai", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://openai.com", w.Header().Get("Location"))
}

func TestRedirectShortURL_Expired(t *testing.T) {
	setupTestDB()

	past := time.Now().Add(-1 * time.Hour)
	database.DB.Create(&models.URL{
		OriginalURL: "https://expired.com",
		ShortCode:   "exp",
		UserID:      1,
		ExpiresAt:   &past,
	})

	r := setupRouter()
	r.GET("/short/:code", RedirectShortURL)

	req, _ := http.NewRequest("GET", "/short/exp", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusGone, w.Code)
}

func TestDeleteURL_Forbidden(t *testing.T) {
	setupTestDB()

	database.DB.Create(&models.URL{
		OriginalURL: "https://secret.com",
		ShortCode:   "sec",
		UserID:      99,
	})

	r := setupRouter()
	r.DELETE("/urls/:id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		DeleteURL(c)
	})

	req, _ := http.NewRequest("DELETE", "/urls/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
