/*
   Data de criação: 27/09/2026
   Estudante: Gustavo Martins França
   Definição do servidor
*/

package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sockets_tcp/Q1/internal/auth"

	"github.com/joho/godotenv"
)

var auth_users = map[string]string{
	"carol": auth.HashSHA512("carol123"),
	"duda":  auth.HashSHA512("duda123"),
}

type Server struct {
	addr     string
	baseDir  string
	listener net.Listener
}

// Server constructor
func NewServer(addr string, baseDir string) *Server {
	return &Server{
		addr:    addr,
		baseDir: baseDir,
	}
}

// Listens to addr and creates goroutines to handle each session
func (s *Server) Run() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)

	if err != nil {
		log.Fatalf("Error listening to %s", s.addr)
	}

	defer s.listener.Close()

	err = os.MkdirAll("logs/", os.ModePerm)
	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatalf("Could not open log file: %s", err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Println("Error accepting connection")
			continue
		}

		session := NewSession(conn, s.baseDir)
		go session.handleConnection()
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading enviroment")
	}

	addr := os.Getenv("SERVER_ADDR")
	baseDir := os.Getenv("BASE_DIR")

	if addr == "" || baseDir == "" {
		log.Fatal("Provide a valid enviroment")
	}

	// dirty dir creation for every user; should be populated with some files manually later
	for user := range auth_users {
		userPath := filepath.Join(baseDir, user)
		err := os.MkdirAll(userPath, os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating directory for user %s", user)
		}
	}

	server := NewServer(addr, baseDir)
	server.Run()
}
