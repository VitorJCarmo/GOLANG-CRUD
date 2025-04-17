package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

type ApiError struct {
	erroMsg string
}
type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle error
			writeJSON(w, r.Response.StatusCode, ApiError{erroMsg: err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func NewAPIServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAcc))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccById))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("Api running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAcc(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAcc(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAcc(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAcc(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAcc(w, r)
	}
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAcc(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		fmt.Println("Erro ao obter todas as contas: ", err.Error())
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]
	id, err := strconv.Atoi(vars)
	if err != nil {
		panic(err)
	}

	acc, err := s.storage.GetAccountByID(id)
	if err != nil {
		fmt.Println("Erro ao obter conta: ", err.Error())
	}
	return writeJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateAcc(w http.ResponseWriter, r *http.Request) error {
	createReq := new(CreateAccountReq)

	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		fmt.Println("Não foi possível decodificar o json: ", err.Error())
		return err
	}

	newAccount := NewAccount(createReq.FirstName, createReq.LastName)

	if id, err := s.storage.CreateAccount(newAccount); err != nil {
		fmt.Println("Não foi possível criar conta: ", err.Error())
		return err
	} else {
		newAccount.ID = id
		fmt.Println("Conta criada", id)
	}

	return writeJSON(w, http.StatusOK, newAccount)
}

func (s *APIServer) handleDeleteAcc(w http.ResponseWriter, r *http.Request) error {
	deleteReq := new(DeleteAccountReq)
	if err := json.NewDecoder(r.Body).Decode(deleteReq); err != nil {
		fmt.Println("Não foi possível decodificar o json: ", err.Error())
		return err
	}
	count, err := s.storage.DeleteAccount(deleteReq.ID)
	if err != nil {
		fmt.Println("Não foi possível deletar conta: ", err.Error())
	}
	return writeJSON(w, http.StatusOK, count)
}

func (s *APIServer) handleUpdateAcc(w http.ResponseWriter, r *http.Request) error {
	createReq := new(UpdateAccountReq)

	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		fmt.Println("Não foi possível decodificar o json: ", err.Error())
		return err
	}

	newAccount := Account{
		ID:        createReq.ID,
		FirstName: createReq.FirstName,
		LastName:  createReq.LastName,
		Number:    createReq.Number,
		Balance:   createReq.Balance,
	}

	if count, err := s.storage.UpdateAccount(newAccount); err != nil {
		fmt.Println("Não foi possível atualizar conta: ", err.Error())
		return err
	} else {
		fmt.Println(count, "Linhas foram atualizadas")
	}

	return writeJSON(w, http.StatusOK, newAccount)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {

		transferReq := new(TransferReq)
		if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
			fmt.Println("Não foi possível decodificar o json: ", err.Error())
			return err
		}
		contaOrigem, erro := s.storage.GetAccountByID(transferReq.ID)
		if erro != nil {
			fmt.Println("Não foi possível obter conta origem: ", erro.Error())
		}

		contaDestino, errd := s.storage.GetAccountByID(transferReq.IdDestino)
		if errd != nil {
			fmt.Println("Não foi possível obter conta destino: ", erro.Error())
		}

		if contaOrigem.Balance < transferReq.Valor {
			fmt.Println("Saldo Insuficiente: ", erro.Error())
			return writeJSON(w, http.StatusOK, transferReq.Valor)
		}

		contaOrigem.Balance = contaOrigem.Balance - transferReq.Valor
		contaDestino.Balance = contaDestino.Balance + transferReq.Valor
		s.storage.UpdateAccount(contaOrigem)
		s.storage.UpdateAccount(contaDestino)
		return writeJSON(w, http.StatusOK, nil)
	}
	return nil
}
