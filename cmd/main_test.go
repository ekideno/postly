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
		Database: config.Database{
			Host:     "localhost",
			User:     "postgres",
			Password: 1234,
			Name:     "postly",
			Port:     5432,
			Sslmode:  "disable",
			Timezone: "UTC",
		},
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := `{"hello":"Postly!"}`
	assert.JSONEq(t, expected, w.Body.String())
}
