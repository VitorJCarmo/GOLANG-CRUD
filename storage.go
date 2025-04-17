package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	GetAccounts() ([]Account, error)
	GetAccountByID(id int) (Account, error)
	DeleteAccount(id int) (int64, error)
	UpdateAccount(Account) (int64, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	//connStr := "user=postgres dbname=postgres password=admin sslmode=disable"
	db, err := sql.Open("postgres", "postgres://postgres:admin@postgres/postgres?sslmode=disable")
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

func (s *PostgresStore) CreateAccount(acc *Account) (int, error) {
	lastInsertId := 0
	query := "insert into account(first_name,last_name,number,balance,created_at) values($1,$2,$3,$4,$5) RETURNING id"
	s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt).Scan(&lastInsertId)
	if lastInsertId == 0 {
		fmt.Println("Erro ao inserir conta")
	}
	return lastInsertId, nil
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

func (s *PostgresStore) GetAccountByID(id int) (Account, error) {
	var acc Account
	query := `select id,first_name,last_name,number,balance,created_at from account where id =$1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return acc, err
	}
	for rows.Next() {
		if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName,
			&acc.Number, &acc.Balance, &acc.CreatedAt); err != nil {
			return acc, err
		}
	}

	return acc, nil
}

func (s *PostgresStore) DeleteAccount(id int) (int64, error) {
	query := "delete from account where id = $1"
	result, err := s.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, err
	}
	count, nil := result.RowsAffected()
	return count, nil
}

func (s *PostgresStore) UpdateAccount(acc Account) (int64, error) {
	query := "update account set first_name = $1, last_name = $2, number = $3, balance = $4  where id = $5"
	result, err := s.db.ExecContext(context.Background(), query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.ID)
	if err != nil {
		return 0, err
	}
	count, nil := result.RowsAffected()
	return count, nil
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
