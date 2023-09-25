package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
	"rsoi-lab-1/internal/model"
	"strconv"
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

func (h *Handler) GetPerson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			err = fmt.Errorf("parse id from query to uint64: %w", err)
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		person, err := h.repository.Get(id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = fmt.Errorf("person not found: %w", err)
				h.WriteError(w, err, http.StatusNotFound)
			} else {
				h.WriteError(w, err, http.StatusInternalServerError)
			}
			return
		}

		err = h.WriteResponse(w, person.ToResponse(), http.StatusOK, nil)
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
					fmt.Sprintf("/api/v1/persons/%d", id),
				},
			},
		)

		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h *Handler) DeletePerson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			err = fmt.Errorf("parse id from query to uint64: %w", err)
			h.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		err = h.repository.Delete(id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = fmt.Errorf("person not found: %w", err)
				h.WriteError(w, err, http.StatusNotFound)
			} else {
				h.WriteError(w, err, http.StatusInternalServerError)
			}
			return
		}

		err = h.WriteResponse(w, nil, http.StatusNoContent, nil)

		if err != nil {
			h.WriteError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h *Handler) UpdatePerson() func(w http.ResponseWriter, r *http.Request) {
	return nil
}
