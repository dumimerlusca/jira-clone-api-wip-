package mock

import (
	"net/http"
)

type MockHandler struct {
	CallsCount int
}

func (h *MockHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	h.CallsCount += 1

	w.WriteHeader(http.StatusOK)
}

func NewMockHandler() *MockHandler {
	return &MockHandler{CallsCount: 0}
}
