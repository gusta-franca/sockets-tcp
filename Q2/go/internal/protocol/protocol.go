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

	// write common request header
	header := []byte{
		byte(req.MessageType),
		byte(req.CommandIdentifier),
		byte(filenameSize),
	}

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

// Formats a response to the binary protocol and sends it to the client
func (p *Protocol) WriteResponse(res *Response) error {
	// write common request header
	header := []byte{
		byte(res.MessageType),
		byte(res.CommandIdentifier),
		byte(res.StatusCode),
	}

	if _, err := p.conn.Write(header); err != nil {
		return err
	}

	// sends header only on error
	if res.StatusCode == StatusError {
		return nil
	}

	// write additional GETFILESLIST fields
	if res.CommandIdentifier == CmdGetFilesList {
		fileCount := make([]byte, 2)
		filesLen := len(res.Files)
		binary.BigEndian.PutUint16(fileCount, uint16(filesLen))

		if _, err := p.conn.Write(fileCount); err != nil {
			return err
		}

		filenameSize := make([]byte, 1)

		for i := range filesLen {
			filename := []byte(res.Files[i])
			filenameLen := len(filename)
			filenameSize[0] = byte(filenameLen)

			if _, err := p.conn.Write(filenameSize); err != nil {
				return err
			}

			if _, err := p.conn.Write(filename); err != nil {
				return err
			}
		}
	}

	// write additional GETFILE fields
	if res.CommandIdentifier == CmdGetFile {
		fileSize := make([]byte, 4)
		fileLen := len(res.File)
		binary.BigEndian.PutUint32(fileSize, uint32(fileLen))

		if _, err := p.conn.Write(fileSize); err != nil {
			return err
		}

		if _, err := p.conn.Write(res.File); err != nil {
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
