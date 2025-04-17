package main

import (
	"math/rand"
	"time"
)

type CreateAccountReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateAccountReq struct {
	ID        int     `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Number    int64   `json:"number"`
	Balance   float64 `json:"balance"`
}

type DeleteAccountReq struct {
	ID int `json:"id"`
}

type TransferReq struct {
	ID        int     `json:"id"`
	IdDestino int     `json:"idDestino"`
	Valor     float64 `json:"valor"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		//ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(100000)),
		CreatedAt: time.Now().UTC(),
	}
}
