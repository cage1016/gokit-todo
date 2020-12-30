// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	"github.com/cage1016/gokit-todo/internal/pkg/responses"
)

type response struct {
	responses.DataRes
	Data []model.Todo `json:"data"`
}

func Test_Todo(t *testing.T) {
	t.Cleanup(func() {
		if err := Truncate(a.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	})

	// add todo
	req, _ := http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"aa","completed":false}`))
	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// add todo
	req, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"bb","completed":false}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// add todo
	req, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"cc","completed":false}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// list todos
	req, _ = http.NewRequest(http.MethodGet, "/items", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	res := response{}
	json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, 3, len(res.Data))

	// edit text
	req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/items/%s", res.Data[0].ID),
		strings.NewReader(`{"completed":false,"text":"dd"}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("status: excpet 200, got %d", w.Code))

	// set completed
	req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/items/%s", res.Data[0].ID),
		strings.NewReader(`{"completed":true}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("status: excpet 200, got %d", w.Code))

	// delete
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/items/%s", res.Data[0].ID), nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, fmt.Sprintf("status: excpet 200, got %d", w.Code))

	// list todos all complete
	req, _ = http.NewRequest(http.MethodGet, "/items", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	res = response{}
	json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, 2, len(res.Data))
}
