/*
   Data de criação: 27/09/2026
   Estudante: Gustavo Martins França
   Definição de uma sessão entre o servidor e um usuário
*/

package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sd/sockets_tcp/internal/protocol"
	"strings"
)

type Session struct {
	conn    net.Conn
	proto   *protocol.Utf8Protocol
	auth    bool
	user    string
	baseDir string
	rootDir string
	currDir string
}

// Session constructor
func NewSession(conn net.Conn, userDir string) *Session {
	return &Session{conn: conn,
		proto: protocol.NewUtf8Protocol(conn), auth: false,
		baseDir: userDir,
	}
}

// Standardizes the client's input string and calls the command handler
func (s *Session) handleConnection() {
	defer s.conn.Close()
	log.Printf("Client connected: %s", s.conn.RemoteAddr())

	for {
		message, err := s.proto.ReadString()

		if err != nil {
			log.Printf("Connection finished: %s", s.conn.RemoteAddr())
			return
		}

		command := strings.TrimSpace(message)

		if command == "" {
			continue
		}

		if !s.HandleCommand(command) {
			log.Printf("Connection finished: %s", s.conn.RemoteAddr())
			return
		}
	}
}

// Handles each command appropriatly
func (s *Session) HandleCommand(command string) bool {
	tokens := strings.SplitN(command, " ", 2)
	cmd := strings.ToUpper(tokens[0])
	args := ""

	// tokens is a comma separated vector, remember to split later
	if len(tokens) > 1 {
		args = strings.TrimSpace(tokens[1])
	}

	if cmd == "EXIT" {
		s.serverLog(command)

		return false
	}

	if cmd == "CONNECT" {
		s.Connect(args)
		s.serverLog(command)

		return true
	}

	if !s.auth {
		s.sendResponse("ERROR")
		return true
	}

	switch cmd {
	case "PWD":
		s.PrintWorkingDirectory()

	case "CHDIR":
		s.ChangeDirectory(args)

	case "GETFILES":
		s.GetFiles()

	case "GETDIRS":
		s.GetDirs()

	default:
		s.sendResponse("ERROR")
	}

	s.serverLog(command)
	return true
}

// Handles CONNECT user,password
func (s *Session) Connect(args string) {
	tokens := strings.Split(args, ",")

	if len(tokens) != 2 {
		s.sendResponse("ERROR")
		return
	}

	user := strings.TrimSpace(tokens[0])
	hash := strings.TrimSpace(tokens[1])

	if !s.validateUser(user, hash) {
		s.sendResponse("ERROR")
		return
	}

	s.auth = true
	s.user = user
	s.rootDir = filepath.Join(s.baseDir, user)
	s.currDir = s.rootDir
	s.sendResponse("SUCCESS")
}

// Handles PWD
func (s *Session) PrintWorkingDirectory() {
	rel, err := filepath.Rel(s.rootDir, s.currDir)

	// trynig to hide the rest of the filesystem
	if err != nil || rel == "." {
		s.sendResponse("/")
	} else {
		s.sendResponse("/" + filepath.ToSlash(rel))
	}
}

// Handles CHDIR
func (s *Session) ChangeDirectory(dir string) {
	if dir == "" {
		s.sendResponse("ERROR")
		return
	}

	if dir == "." {
		s.sendResponse("SUCCESS")
		return
	}

	var target string

	if filepath.IsAbs(dir) {
		target = filepath.Join(s.rootDir, filepath.Clean(dir))
	} else {
		target = filepath.Join(s.currDir, dir)
	}

	target = filepath.Clean(target)

	// Rejects command if the required path doesn't begin with {user}/
	if !strings.HasPrefix(target, s.rootDir) {
		s.sendResponse("ERROR")
		return
	}

	info, err := os.Stat(target)

	if err != nil || !info.IsDir() {
		s.sendResponse("ERROR")
		return
	}

	s.currDir = target
	s.sendResponse("SUCCESS")
}

// Handles GETFILES
func (s *Session) GetFiles() {
	c, err := os.ReadDir(s.currDir)

	if err != nil {
		s.sendResponse("ERROR")
		return
	}

	var files []string
	for _, entry := range c {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	s.sendResponse(strings.Join(files, "\n"))
}

// Handles GETDIRS
func (s *Session) GetDirs() {
	c, err := os.ReadDir(s.currDir)

	if err != nil {
		s.sendResponse("ERROR")
		return
	}

	var dirs []string
	for _, entry := range c {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	s.sendResponse(strings.Join(dirs, "\n"))
}

// Formats and send a command's response to the client
func (s *Session) sendResponse(response string) {
	err := s.proto.WriteString(strings.ToValidUTF8(response, ""))

	if err != nil {
		log.Printf("Error sending response: %s", err)
	}
}

// Validates if an user is authenticated
func (s *Session) validateUser(user string, hash string) bool {
	realHash, ok := auth_users[user] // ok means "exists"

	return ok && strings.EqualFold(hash, realHash)
}

// Logs every executed command
func (s *Session) serverLog(command string) {
	log.Printf("%s executed %s\n", s.conn.RemoteAddr(), command)
}
