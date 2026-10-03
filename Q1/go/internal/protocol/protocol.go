/*
   Data de criação: 02/10/2026
   Estudante: Gustavo Martins França
   Definição do protocolo comum entre linguagens
*/

package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

type Utf8Protocol struct {
	conn net.Conn
}

// Utf8Protocol constructor
func NewUtf8Protocol(conn net.Conn) *Utf8Protocol {
	return &Utf8Protocol{
		conn: conn,
	}
}

// Reads a big-endian header with 4 bytes, then reads the bytes informed in the header
func (p *Utf8Protocol) ReadString() (string, error) {
	header := make([]byte, 4)

	// reads header
	if _, err := io.ReadFull(p.conn, header); err != nil {
		return "", err
	}

	length := binary.BigEndian.Uint32(header)
	data := make([]byte, length)

	// reads data
	if _, err := io.ReadFull(p.conn, data); err != nil {
		return "", err
	}

	return string(data), nil
}

// Converts a string to UTF-8, writes a big endian header with 4 bytes, then writes the body
func (p *Utf8Protocol) WriteString(message string) error {
	data := []byte(message)
	length := uint32(len(data))

	header := make([]byte, 4)

	binary.BigEndian.PutUint32(header, length)

	// writes header
	if _, err := p.conn.Write(header); err != nil {
		return err
	}

	// writes data
	_, err := p.conn.Write(data)
	return err
}
