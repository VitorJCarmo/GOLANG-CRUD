package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAcc(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		log.Fatal("Erro ao obters todas as contas:", err.Error())
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]
	//account := NewAccount("vitor", "carmo")
	return writeJSON(w, http.StatusOK, vars)
}

func (s *APIServer) handleCreateAcc(w http.ResponseWriter, r *http.Request) error {
	createReq := new(CreateAccountReq)
	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		log.Fatal("Não foi possível decodificar o json ", err.Error())
		return err
	}
	newAccount := NewAccount(createReq.FirstName, createReq.LastName)
	if err := s.storage.CreateAccount(newAccount); err != nil {
		log.Fatal("Não foi possível criar nova conta ", err.Error())
		return err
	}
	return writeJSON(w, http.StatusOK, newAccount)
}

func (s *APIServer) handleDeleteAcc(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
