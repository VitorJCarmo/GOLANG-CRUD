package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccounts() ([]Account, error)
	GetAccountByID(id int) error
	DeleteAccount(id int) error
	UpdateAccount(*Account) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=admin sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := "insert into account(id,first_name,last_name,number,balance,created_at) values($1,$2,$3,$4,$5,$6)"
	_, err := s.db.ExecContext(context.Background(), query, acc.ID, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccounts() ([]Account, error) {
	query := `select id,first_name,last_name,number,balance,created_at from account`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName,
			&acc.Number, &acc.Balance, &acc.CreatedAt); err != nil {
			return accounts, err
		}
		accounts = append(accounts, acc)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(int) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}
