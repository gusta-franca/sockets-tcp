/*
   Data de criação: 03/10/2026
   Estudante: Gustavo Martins França
   Definição do protocolo comum entre linguagens
*/

package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Protocol struct {
	conn net.Conn
}

// Protocol constructor
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		conn: conn,
	}
}

// Formats a request to the binary protocol and sends it to the server
func (p *Protocol) WriteRequest(req *Request) error {
	filename := []byte(req.Filename)
	filenameSize := len(filename)

	if filenameSize > 255 {
		return fmt.Errorf("Filename size exceeds 255 bytes limit")
	}

	header := []byte{
		byte(req.MessageType),
		byte(req.CommandIdentifier),
		byte(filenameSize),
	}

	// write common request header
	if _, err := p.conn.Write(header); err != nil {
		return err
	}

	// write filename if not empty
	if filenameSize > 0 {
		if _, err := p.conn.Write(filename); err != nil {
			return err
		}
	}

	// write additional ADDFILE fields
	if req.CommandIdentifier == CmdAddFile {
		if req.FileSize == nil {
			return fmt.Errorf("FileSize is nil")
		}

		fileSize := make([]byte, 4)
		binary.BigEndian.PutUint32(fileSize, *req.FileSize)

		if _, err := p.conn.Write(fileSize); err != nil {
			return err
		}

		if _, err := p.conn.Write(req.File); err != nil {
			return err
		}
	}

	return nil
}

// Reads a protocol message and returns a Request struct with the appropriate fields
func (p *Protocol) ReadRequest() (*Request, error) {
	header := make([]byte, 3)

	// read header
	if _, err := io.ReadFull(p.conn, header); err != nil {
		return nil, err
	}

	req := &Request{
		MessageType:       Message(header[0]),
		CommandIdentifier: Command(header[1]),
	}

	// read filename if filenameSize > 0
	filenameSize := header[2]

	if filenameSize > 0 {
		filename := make([]byte, filenameSize)

		if _, err := io.ReadFull(p.conn, filename); err != nil {
			return nil, err
		}

		req.Filename = string(filename)
	}

	// read ADDFILE fields
	if req.CommandIdentifier == CmdAddFile {
		fileSize := make([]byte, 4)

		if _, err := io.ReadFull(p.conn, fileSize); err != nil {
			return nil, err
		}

		// needed because FileSize is a pointer, it's not possible to do *req.FileSize = binary.BigEndian.Uint32(fileSize)
		s := binary.BigEndian.Uint32(fileSize)
		req.FileSize = &s

		file := make([]byte, *req.FileSize)

		if _, err := io.ReadFull(p.conn, file); err != nil {
			return nil, err
		}

		req.File = file
	}

	return req, nil
}
