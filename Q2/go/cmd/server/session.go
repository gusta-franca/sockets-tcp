package server

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Session struct {
	conn     net.Conn
	reader   *bufio.Reader
	auth     bool
	user     string
	curr_dir string
}

// onde colocar o diretório? criar um em /temp/user para ser o padrão?
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn, reader: bufio.NewReader(conn), auth: false}
}

// HandleConnection; passa o comando para frente com a string padronizada
func (s *Session) handleConnection() {
	defer s.conn.Close()
	log.Printf("Cliente conectado - %s", s.conn.RemoteAddr())

	for {
		message, err := s.reader.ReadString('\n')

		if err != nil {
			log.Printf("Erro durante ReadString - %s; conexão encerrada: %v", s.conn.RemoteAddr(), err)

			return
		}

		command := strings.TrimSpace(message)

		if command == "" {
			continue
		}

		if !s.HandleCommand(command) {
			log.Printf("Conexão encerrada - %s", s.conn.RemoteAddr())

			return
		}
	}
}

// HandleCommand; switch com casos == comandos na especificação
func (s *Session) HandleCommand(command string) bool {
	tokens := strings.Split(command, " ")
	cmd := strings.ToUpper(tokens[0])

	switch cmd {

	case "CONNECT":
		s.Connect(tokens[1], tokens[2])

		return true

	case "PWD":
		s.PrintWorkingDirectory()

		return true

	case "CHDIR":
		s.ChangeDirectory()

		return true

	case "GETFILES":
		s.GetFiles()

		return true

	case "GETDIRS":
		s.GetDirs()

		return true

	case "EXIT":
		return false
	}
}
