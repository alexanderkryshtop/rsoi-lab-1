package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rsoi-lab-1/internal/model"
)

func (h *Handler) GetAllPersons() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(r.Body)

		repository := h.repository
		persons, err := repository.GetAll()

		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		personResponses := make([]model.PersonResponse, 0)
		for _, person := range persons {
			personResponse := person.ToResponse()
			personResponses = append(personResponses, personResponse)
		}

		err = h.WriteResponse(w, personResponses, http.StatusOK, nil)
		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h *Handler) CreatePerson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := io.ReadAll(r.Body)

		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(r.Body)

		if err != nil {
			err = fmt.Errorf("http request body read: %w", err)
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		personRequest := new(model.PersonRequest)
		err = json.Unmarshal(requestBody, personRequest)
		if err != nil {
			err = fmt.Errorf("json unmarshal: %w", err)
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		person := model.FromRequest(personRequest)
		id, err := h.repository.Create(person)
		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		err = h.WriteResponse(w, http.NoBody, http.StatusCreated,
			map[string][]string{
				"Located": {
					fmt.Sprintf("/api/v1/persons/{%d}", id),
				},
			},
		)

		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
		}
	}
}
