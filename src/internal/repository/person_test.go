package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsoi-lab-1/internal/model"
	"rsoi-lab-1/internal/repository/mocks"
	"testing"
)

//go:generate mockery --name Repository

func ptr[T any](v T) *T {
	return &v
}

func TestPersonRepository_Get(t *testing.T) {
	repository := mocks.Repository{}

	expectedPerson := &model.Person{
		ID:      1,
		Name:    "Alice",
		Age:     ptr[int32](20),
		Address: ptr("Alice address"),
		Work:    ptr("Alice work"),
	}
	repository.On("Get", expectedPerson.ID).Return(expectedPerson, nil)

	person, err := repository.Get(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPerson, person)
}

func TestPersonRepository_GetAll(t *testing.T) {
	repository := mocks.Repository{}

	expectedPersons := []model.Person{
		{
			ID:      1,
			Name:    "Alice",
			Age:     ptr[int32](20),
			Address: ptr("Alice address"),
			Work:    ptr("Alice work"),
		},
		{
			ID:      2,
			Name:    "Bob",
			Age:     ptr[int32](22),
			Address: ptr("Bob address"),
			Work:    ptr("Bob work"),
		},
	}
	repository.On("GetAll").Return(expectedPersons, nil)

	persons, err := repository.GetAll()

	assert.NoError(t, err)
	assert.EqualValues(t, expectedPersons, persons)
}

func TestPersonRepository_Create(t *testing.T) {
	repository := mocks.Repository{}

	expectedPerson := &model.Person{
		ID:      1,
		Name:    "Alice",
		Age:     ptr[int32](20),
		Address: ptr("Alice address"),
		Work:    ptr("Alice work"),
	}
	repository.On("Create", expectedPerson).Return(expectedPerson.ID, nil)

	personID, err := repository.Create(expectedPerson)

	assert.NoError(t, err)
	assert.EqualValues(t, expectedPerson.ID, personID)
}

func TestPersonRepository_Update(t *testing.T) {
	repository := mocks.Repository{}

	person := &model.Person{
		ID:   1,
		Name: "Alice",
	}
	updatedPerson := &model.Person{
		ID:      1,
		Name:    "Alice",
		Age:     ptr[int32](20),
		Address: ptr("Alice address"),
		Work:    ptr("Alice work"),
	}

	repository.On("Update", person).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.Person)
		arg.Age = updatedPerson.Age
		arg.Address = updatedPerson.Address
		arg.Work = updatedPerson.Work
	}).Return(nil)
	err := repository.Update(person)

	assert.NoError(t, err)
	assert.EqualValues(t, updatedPerson, person)
}

func TestPersonRepository_Delete(t *testing.T) {
	repository := mocks.Repository{}

	person := &model.Person{
		ID:      1,
		Name:    "Alice",
		Age:     ptr[int32](20),
		Address: ptr("Alice address"),
		Work:    ptr("Alice work"),
	}
	repository.On("Delete", person.ID).Return(nil)

	err := repository.Delete(person.ID)

	assert.NoError(t, err)
}
