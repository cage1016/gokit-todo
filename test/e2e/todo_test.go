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

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/pkg/responses"
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
	req, _ := http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"aa"}`))
	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// add todo
	req, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"bb"}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// add todo
	req, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"text":"cc"}`))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, fmt.Sprintf("status: excpet 201, got %d", w.Code))

	// list todos
	req, _ = http.NewRequest(http.MethodGet, "/items?filter=all", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	res := response{}
	json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, 3, len(res.Data))

	// update complete
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/items/%s", res.Data[0].ID),
		strings.NewReader(fmt.Sprintf(`{"id":"%s","text":"dd"}`, res.Data[0].ID)))
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("status: excpet 200, got %d", w.Code))

	// complete
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/items/completes/%s", res.Data[0].ID), nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, fmt.Sprintf("status: excpet 204, got %d", w.Code))

	// complete all
	req, _ = http.NewRequest(http.MethodPost, "/items/completes", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, fmt.Sprintf("status: excpet 204, got %d", w.Code))

	// list todos all complete
	req, _ = http.NewRequest(http.MethodGet, "/items?filter=complete", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	res = response{}
	json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, 3, len(res.Data))
	assert.Equal(t, true, res.Data[0].Complete)
	assert.Equal(t, true, res.Data[1].Complete)
	assert.Equal(t, true, res.Data[2].Complete)

	// clear complete
	req, _ = http.NewRequest(http.MethodPost, "/items/clear", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, fmt.Sprintf("status: excpet 204, got %d", w.Code))

	// list all todos
	req, _ = http.NewRequest(http.MethodGet, "/items?filter=all", nil)
	w = httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	res = response{}
	json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, 0, len(res.Data))
}
