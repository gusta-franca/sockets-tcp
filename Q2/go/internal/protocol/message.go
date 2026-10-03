/*
   Data de criação: 03/10/2026
   Estudante: Gustavo Martins França
   Definição das mensagens do protocolo
*/

package protocol

type Message byte
type Command byte
type Status byte

// no enum in Go
const (
	MessageRequest  Message = 0x01
	MessageResponse Message = 0x02
)

const (
	CmdAddFile      Command = 0x01
	CmdDelete       Command = 0x02
	CmdGetFilesList Command = 0x03
	CmdGetFile      Command = 0x04
)

const (
	StatusSuccess Status = 0x01
	StatusError   Status = 0x02
)

type Request struct {
	MessageType       Message
	CommandIdentifier Command
	Filename          string
	// ADDFILE
	FileSize *uint32 // pointer so an unset FileSize is always nil
	File     []byte  // 1-2^32 bytes
}

type Response struct {
	MessageType       Message
	CommandIdentifier Command
	StatusCode        Status
	// GETFILESLIST
	Files []string
	// GETFILE
	FileSize *uint32
	File     []byte
}
