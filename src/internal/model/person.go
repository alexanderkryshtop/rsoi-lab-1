package model

import "fmt"

type Person struct {
	ID      int32
	Name    string
	Age     *int32
	Address *string
	Work    *string
}

type PersonRequest struct {
	Name    string  `json:"name"`
	Age     *int32  `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

type PersonResponse struct {
	ID      int32   `json:"id"`
	Name    string  `json:"name"`
	Age     *int32  `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

func FromRequest(r *PersonRequest) (*Person, error) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}
	return &Person{
		Name:    r.Name,
		Age:     r.Age,
		Address: r.Address,
		Work:    r.Work,
	}, nil
}

func (p *Person) ToResponse() PersonResponse {
	return PersonResponse{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}

type PersonValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (e PersonValidationError) Error() string {
	return fmt.Sprintf("%s %v", e.Message, e.Errors)
}

func (p *PersonRequest) Validate() error {
	e := &PersonValidationError{
		Message: "Fields validation failed",
		Errors:  map[string]string{},
	}
	if p.Age != nil && *(p.Age) < 0 {
		e.Errors["age"] = fmt.Sprintf("must be greater or equal 0, now it is %d", *p.Age)
	}
	if p.Work != nil && len(*p.Work) == 0 {
		e.Errors["work"] = fmt.Sprintf("must not be empty string")
	}
	if p.Address != nil && len(*p.Address) == 0 {
		e.Errors["work"] = fmt.Sprintf("must not be empty string")
	}

	if len(e.Errors) > 0 {
		return e
	}
	return nil
}
