package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	todoist "github.com/niskhakov/gotodoist"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

type AuthorizationServer struct {
	server          *http.Server
	todoistClient   *todoist.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
	logger          *log.Logger
}

func NewAuthorizationServer(todoistClient *todoist.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{todoistClient: todoistClient, tokenRepository: tokenRepository, redirectURL: redirectURL, logger: log.New(os.Stdout, "AUTHSERV: ", log.LstdFlags)}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":3001",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	errorParam := r.URL.Query().Get("error")
	if errorParam != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: " + errorParam))
		return
	}

	chatIDParam := r.URL.Query().Get("state")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Param state was not provided"))
		return
	}

	codeParam := r.URL.Query().Get("code")
	if codeParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Code was not specified by Todoist app"))
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State param authenticity issue"))
		return
	}

	err = s.tokenRepository.Save(chatID, codeParam, repository.CodeTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	accessToken, err := s.todoistClient.GetAccessToken(r.Context(), codeParam)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err.Error())
		return
	}

	err = s.tokenRepository.Save(chatID, accessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	s.logger.Printf("chat id: %d\ncode token: %s\naccess_token: %s\n", chatID, codeParam, accessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
