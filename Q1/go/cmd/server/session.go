package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Session struct {
	conn    net.Conn
	reader  *bufio.Reader
	auth    bool
	user    string
	baseDir string
	rootDir string
	currDir string
}

// onde colocar o diretório? criar um em /temp/user para ser o padrão?
func NewSession(conn net.Conn, userDir string) *Session {
	return &Session{conn: conn,
		reader:  bufio.NewReader(conn),
		auth:    false,
		baseDir: userDir,
	}
}

// HandleConnection; passa o comando para frente com a string padronizada
func (s *Session) handleConnection() {
	defer s.conn.Close()
	log.Printf("Client connected: %s", s.conn.RemoteAddr())

	for {
		message, err := s.reader.ReadString('\n')

		if err != nil {
			s.sendResponse("ERROR")
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

// HandleCommand; switch com casos == comandos na especificação
func (s *Session) HandleCommand(command string) bool {
	tokens := strings.SplitN(command, " ", 2)
	cmd := strings.ToUpper(tokens[0])
	args := ""

	// tokens é um vetor separado por vírgula, lembrar de remover vírgulas se necessário
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

func (s *Session) PrintWorkingDirectory() {
	rel, err := filepath.Rel(s.rootDir, s.currDir)

	// trynig to hide the rest of the filesystem
	if err != nil || rel == "." {
		s.sendResponse("/")
	} else {
		s.sendResponse("/" + filepath.ToSlash(rel))
	}
}

func (s *Session) ChangeDirectory(dir string) {
	if dir == "" {
		s.sendResponse("ERROR")
		return
	}

	if dir == "." {
		return
	}

	var target string

	if filepath.IsAbs(dir) {
		target = filepath.Join(s.rootDir, filepath.Clean(dir))
	} else {
		target = filepath.Join(s.currDir, dir)
	}

	target = filepath.Clean(target)

	// rejects command if the required path doesn't begin with {user}/
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

	var response strings.Builder

	response.WriteString(strconv.Itoa(len(files)))
	response.WriteString("\n")

	for _, f := range files {
		response.WriteString(f)
		response.WriteString("\n")
	}

	s.sendResponse(strings.TrimSuffix(response.String(), "\n"))
}

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

	var response strings.Builder

	response.WriteString(strconv.Itoa(len(dirs)))
	response.WriteString("\n")

	for _, d := range dirs {
		response.WriteString(d)
		response.WriteString("\n")
	}

	s.sendResponse(strings.TrimSuffix(response.String(), "\n"))
}

func (s *Session) sendResponse(response string) {
	_, err := fmt.Fprintf(s.conn, "%s\n", strings.ToValidUTF8(response, ""))

	if err != nil {
		log.Printf("Error sending response")
	}
}

func (s *Session) validateUser(user string, hash string) bool {
	realHash, ok := auth_users[user] // ok means "exists"

	return ok && strings.EqualFold(hash, realHash)
}

func (s *Session) serverLog(command string) {
	log.Printf("%s executed %s\n", s.conn.RemoteAddr(), command)
}
