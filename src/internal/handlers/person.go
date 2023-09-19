package handlers

import (
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

		h.WriteResponse(w, personResponses)
	}
}
