package model

type Person struct {
	ID      uint64
	Name    string
	Age     *uint64
	Address *string
	Work    *string
}

type PersonRequest struct {
	Name    string  `json:"name"`
	Age     *uint64 `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

type PersonResponse struct {
	ID      uint64  `json:"id"`
	Name    string  `json:"name"`
	Age     *uint64 `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

func FromRequest(r *PersonRequest) *Person {
	return &Person{
		Name:    r.Name,
		Age:     r.Age,
		Address: r.Address,
		Work:    r.Work,
	}
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

func (p *PersonRequest) Validate() {
	return
}
