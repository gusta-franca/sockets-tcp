/*
   Data de criação: 03/10/2026
   Estudante: Gustavo Martins França
   Definição de uma sessão entre o servidor e um usuário
*/

package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sockets_tcp/Q2/internal/protocol"
	"strings"
)

type Session struct {
	conn    net.Conn
	proto   *protocol.Protocol
	baseDir string
}

// Session constructor
func NewSession(conn net.Conn, baseDir string) *Session {
	return &Session{
		conn:    conn,
		proto:   protocol.NewProtocol(conn),
		baseDir: baseDir,
	}
}

// Standardizes the client's input string and calls the command handler
func (s *Session) handleConnection() {
	defer s.conn.Close()
	log.Printf("Client connected: %s", s.conn.RemoteAddr())

	for {
		req, err := s.proto.ReadRequest()

		if err != nil {
			log.Printf("Connection finished: %s", s.conn.RemoteAddr())
			return
		}

		if !s.HandleRequest(req) {
			log.Printf("Connection finished: %s", s.conn.RemoteAddr())
			return
		}
	}
}

// Handles each request appropriatly
func (s *Session) HandleRequest(req *protocol.Request) bool {
	switch req.CommandIdentifier {
	case protocol.CmdAddFile:
		s.AddFile(req)

	case protocol.CmdDelete:
		s.Delete(req)

	case protocol.CmdGetFilesList:
		s.GetFilesList()

	case protocol.CmdGetFile:
		s.GetFile(req)

	default:
		log.Printf("Unknow command received from %s: 0x%X", s.conn.RemoteAddr(), byte(req.CommandIdentifier))
		s.sendResponse(req.CommandIdentifier, protocol.StatusError, nil, nil)
		return false
	}

	s.serverLog(req.CommandIdentifier)
	return true
}

// Handles ADDFILE
func (s *Session) AddFile(req *protocol.Request) {
	filePath, ok := s.resolvePath(req.Filename)

	if !ok || req.FileSize == nil || req.File == nil {
		s.sendResponse(protocol.CmdAddFile, protocol.StatusError, nil, nil)
		return
	}

	// Ensure slice length matches expected payload size
	// if uint32(len(req.File)) != *req.FileSize {
	// 	s.sendResponse(protocol.CmdAddFile, protocol.StatusError, nil, nil)
	// 	return
	// }

	// write file to disk
	err := os.WriteFile(filePath, req.File, os.ModePerm)

	if err != nil {
		log.Printf("Error adding file %s: %v", req.Filename, err)
		s.sendResponse(protocol.CmdAddFile, protocol.StatusError, nil, nil)
		return
	}

	s.sendResponse(protocol.CmdAddFile, protocol.StatusSuccess, nil, nil)
}

// Handles DELETE
func (s *Session) Delete(req *protocol.Request) {
	filePath, ok := s.resolvePath(req.Filename)

	if !ok {
		s.sendResponse(protocol.CmdDelete, protocol.StatusError, nil, nil)
		return
	}

	err := os.Remove(filePath)

	if err != nil {
		log.Printf("Error deleting file %s: %v", req.Filename, err)
		s.sendResponse(protocol.CmdDelete, protocol.StatusError, nil, nil)
		return
	}

	s.sendResponse(protocol.CmdDelete, protocol.StatusSuccess, nil, nil)
}

// Handles GETFILESLIST
func (s *Session) GetFilesList() {
	c, err := os.ReadDir(s.baseDir)

	if err != nil {
		log.Printf("Error reading directory %s: %v", s.baseDir, err)
		s.sendResponse(protocol.CmdGetFilesList, protocol.StatusError, nil, nil)
		return
	}

	var files []string
	for _, entry := range c {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	s.sendResponse(protocol.CmdGetFilesList, protocol.StatusSuccess, files, nil)
}

// Handles GETFILE
func (s *Session) GetFile(req *protocol.Request) {
	filePath, ok := s.resolvePath(req.Filename)

	if !ok {
		s.sendResponse(protocol.CmdGetFile, protocol.StatusError, nil, nil)
		return
	}

	file, err := os.ReadFile(filePath)

	if err != nil {
		log.Printf("Error reading file %s: %v", req.Filename, err)
		s.sendResponse(protocol.CmdGetFile, protocol.StatusError, nil, nil)
		return
	}

	s.sendResponse(protocol.CmdGetFile, protocol.StatusSuccess, nil, file)
}

// Formats and sends binary protocol response to the client
func (s *Session) sendResponse(command protocol.Command, status protocol.Status, files []string, file []byte) {
	var fileSize *uint32

	if file != nil {
		s := uint32(len(file))
		fileSize = &s
	}

	res := &protocol.Response{
		MessageType:       protocol.MessageResponse,
		CommandIdentifier: command,
		StatusCode:        status,
		Files:             files,
		FileSize:          fileSize,
		File:              file,
	}

	if err := s.proto.WriteResponse(res); err != nil {
		log.Printf("Error sending response to %s: %v", s.conn.RemoteAddr(), err)
	}
}

// Logs every executed command
func (s *Session) serverLog(cmd protocol.Command) {
	log.Printf("%s executed command 0x%02X\n", s.conn.RemoteAddr(), byte(cmd))
}

// Helper to avoid adding/getting files from outside the enviroment
func (s *Session) resolvePath(filename string) (string, bool) {
	if filename == "" {
		return "", false
	}

	target := filepath.Clean(filepath.Join(s.baseDir, filename))

	rel, err := filepath.Rel(s.baseDir, target)

	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		log.Printf("Directory traversal attempted by %s: %s", s.conn.RemoteAddr(), filename)
		return "", false
	}

	return target, true
}
