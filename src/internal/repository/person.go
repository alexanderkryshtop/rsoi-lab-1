package repository

import (
	"context"
	"database/sql"
	"rsoi-lab-1/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(model *model.Person) (uint64, error)
	GetAll() ([]model.Person, error)
	Get(id uint64) (model.Person, error)
	Update(model *model.Person) error
	Delete(id uint64) error
}

type PersonRepository struct {
	dbPool *pgxpool.Pool
}

func NewPersonRepository(dbPool *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{
		dbPool: dbPool,
	}
}

func (p *PersonRepository) Create(model *model.Person) (uint64, error) {
	var personID uint64

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

func (p *PersonRepository) Get(id uint64) (model.Person, error) {
	var person model.Person

	err := p.dbPool.QueryRow(context.Background(),
		"SELECT id, name, age, address, work FROM persons WHERE id=$1",
		id).Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Person{}, nil // Record not found
		}
		return model.Person{}, err
	}

	return person, nil
}

func (p *PersonRepository) Update(person *model.Person) error {
	_, err := p.dbPool.Exec(context.Background(),
		"UPDATE persons SET name=$1, age=$2, address=$3, work=$4 WHERE id=$5",
		person.Name, person.Age, person.Address, person.Work, person.ID)
	return err
}

func (p *PersonRepository) Delete(id uint64) error {
	_, err := p.dbPool.Exec(context.Background(),
		"DELETE FROM persons WHERE id=$1",
		id)
	return err
}
