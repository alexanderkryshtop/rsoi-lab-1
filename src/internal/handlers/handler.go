package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rsoi-lab-1/internal/repository"
)

const (
	ContentTypeKey      = "Content-Type"
	TypeApplicationJSON = "application/json"
)

type Handler struct {
	repository repository.Repository
}

func NewHandler(repository repository.Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) WriteResponse(
	w http.ResponseWriter,
	response interface{},
	statusCode int,
	additionalHeaders http.Header,
) error {
	w.Header().Set(ContentTypeKey, TypeApplicationJSON)
	for key, values := range additionalHeaders {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(statusCode)

	if response == nil {
		return nil
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		err = fmt.Errorf("json marshal: %w", err)
		return err
	}

	_, err = w.Write(responseBytes)
	if err != nil {
		err = fmt.Errorf("http response write: %w", err)
		return err
	}

	return nil
}

func (h *Handler) WriteError(w http.ResponseWriter, err error, code int) {
	type errorResponse struct {
		Message string `json:"message"`
	}
	response := new(errorResponse)

	if err != nil {
		response.Message = err.Error()
	}

	w.Header().Set(ContentTypeKey, TypeApplicationJSON)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) WriteErrorStruct(w http.ResponseWriter, err error, code int) {
	w.Header().Set(ContentTypeKey, TypeApplicationJSON)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(err)
}
