package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekideno/postly/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRootEndpoint(t *testing.T) {
	r := setupRouter(&config.Config{
		DB_HOST:     "localhost",
		DB_USER:     "postgres",
		DB_PASSWORD: 1234,
		DB_NAME:     "postly",
		DB_PORT:     5432,
		DB_SSLMODE:  "disable",
		DB_TIMEZONE: "UTC",
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := `{"hello":"Postly!"}`
	assert.JSONEq(t, expected, w.Body.String())
}
