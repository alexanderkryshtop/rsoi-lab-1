package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"rsoi-lab-1/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(model *model.Person) (int32, error)
	GetAll() ([]model.Person, error)
	Get(id int32) (*model.Person, error)
	Update(model *model.Person) error
	Delete(id int32) error
}

type PersonRepository struct {
	dbPool *pgxpool.Pool
}

func NewPersonRepository(dbPool *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{
		dbPool: dbPool,
	}
}

func (p *PersonRepository) Create(model *model.Person) (int32, error) {
	var personID int32

	err := p.dbPool.QueryRow(context.Background(),
		"INSERT INTO tb_persons (name, age, address, work) VALUES($1, $2, $3, $4) RETURNING id",
		model.Name, model.Age, model.Address, model.Work,
	).Scan(&personID)

	return personID, err
}

func (p *PersonRepository) GetAll() ([]model.Person, error) {
	var people []model.Person

	rows, err := p.dbPool.Query(
		context.Background(),
		"SELECT id, name, age, address, work FROM tb_persons",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work); err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

func (p *PersonRepository) Get(id int32) (*model.Person, error) {
	person := new(model.Person)

	err := p.dbPool.QueryRow(context.Background(),
		"SELECT id, name, age, address, work FROM tb_persons WHERE id=$1",
		id).Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return person, nil
}

func (p *PersonRepository) Update(person *model.Person) error {
	var personID int32
	err := p.dbPool.QueryRow(context.Background(),
		"UPDATE tb_persons SET name=$1, age=$2, address=$3, work=$4 WHERE id=$5 RETURNING id",
		person.Name, person.Age, person.Address, person.Work, person.ID).Scan(&personID)
	return err
}

func (p *PersonRepository) Delete(id int32) error {
	_, err := p.dbPool.Exec(context.Background(),
		"DELETE FROM tb_persons WHERE id=$1",
		id)
	return err
}
