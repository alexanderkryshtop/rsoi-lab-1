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

func (p *Person) FromRequest(r *PersonRequest) {
	p.Name = r.Name
	p.Age = r.Age
	p.Address = r.Address
	p.Work = r.Work
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
