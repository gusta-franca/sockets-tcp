/*
   Data de criação: 03/10/2026
   Estudante: Gustavo Martins França
   Definição do servidor
*/

package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	addr     string
	listener net.Listener
	baseDir  string
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

	err = os.MkdirAll(baseDir, os.ModePerm)

	if err != nil {
		log.Fatal("Error creating directory for clients")

	}

	server := NewServer(addr, baseDir)
	server.Run()
}
