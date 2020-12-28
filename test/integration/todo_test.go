// +build integration

package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func Test_Add_Todo(t *testing.T) {
	t.Cleanup(func() {
		if err := Truncate(a.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	})

	body := strings.NewReader(`{"text":"aa"}`)
	req, err := http.NewRequest(http.MethodPost, "/items", body)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 200, got %d", w.Code))
}
