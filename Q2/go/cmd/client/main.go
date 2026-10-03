/*
   Data de criação: 03/10/2026
   Estudante: Gustavo Martins França
   Definição do cliente
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sockets_tcp/Q2/internal/protocol"
	"strings"

	"github.com/joho/godotenv"
)

type Client struct {
	addr        string
	downloadDir string
	conn        net.Conn
	proto       *protocol.Protocol
	inputReader *bufio.Reader
}

// Client constructor
func NewClient(addr string, downloadDir string) *Client {
	return &Client{
		addr:        addr,
		downloadDir: downloadDir,
	}
}

// Connects a client to a server
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.addr)

	if err != nil {
		return err
	}

	c.inputReader = bufio.NewReader(os.Stdin)
	c.proto = protocol.NewProtocol(c.conn)

	fmt.Printf("Connected on %s\n", c.addr)
	return nil
}

// Connects to addr and handles sending/receiving to/from the server
func (c *Client) Run() {
	err := c.Connect()

	if err != nil {
		log.Fatal("Error connecting to server")
	}

	defer c.conn.Close()

	for {
		fmt.Print(">> ")

		str, err := c.inputReader.ReadString('\n')

		if err != nil {
			log.Fatalf("Error reading input: %s", err)
		}

		tokens := strings.Fields(str)

		if len(tokens) == 0 {
			continue
		}

		command := strings.ToUpper(tokens[0])

		if command == "" {
			continue
		}

		// maybe log on server too?
		if command == "EXIT" {
			fmt.Printf("Disconnected from %s\n", c.addr)
			return
		}

		req := c.buildRequest(command, tokens)

		if req == nil {
			continue
		}

		if err := c.proto.WriteRequest(req); err != nil {
			log.Printf("Error sending request: %v", err)
			break
		}

		res, err := c.proto.ReadResponse()

		if err != nil {
			log.Printf("Error reading response: %v", err)
			break
		}

		c.handleResponse(req, res)
	}
}

// Returns a structured request following the binary protocol specification
func (c *Client) buildRequest(command string, tokens []string) *protocol.Request {
	req := &protocol.Request{
		MessageType: protocol.MessageRequest,
	}

	switch command {
	case "ADDFILE":
		if len(tokens) < 3 {
			fmt.Println("Usage: ADDFILE <local_path> <filename>")
			return nil
		}

		localPath := tokens[1]
		filename := tokens[2]
		file, err := os.ReadFile(localPath)

		if err != nil {
			fmt.Printf("Error reading local file: %v\n", err)
			return nil
		}

		size := uint32(len(file))
		req.CommandIdentifier = protocol.CmdAddFile
		req.Filename = filename
		req.FileSize = &size
		req.File = file

	case "DELETE":
		if len(tokens) < 2 {
			fmt.Println("Usage: DELETE <filename>")
			return nil
		}

		req.CommandIdentifier = protocol.CmdDelete
		req.Filename = tokens[1]

	case "GETFILESLIST":
		req.CommandIdentifier = protocol.CmdGetFilesList

	case "GETFILE":
		if len(tokens) < 2 {
			fmt.Println("Usage: GETFILE <filename>")
			return nil
		}

		req.CommandIdentifier = protocol.CmdGetFile
		req.Filename = tokens[1]

	default:
		fmt.Printf("Unknown command: %s\n", command)
		return nil
	}

	return req
}

// Decodes the server response and prints it
func (c *Client) handleResponse(req *protocol.Request, res *protocol.Response) {
	if res.StatusCode == protocol.StatusError {
		fmt.Printf("Server returned ERROR for command 0x%X\n", req.CommandIdentifier)
		return
	}

	switch res.CommandIdentifier {
	case protocol.CmdAddFile:
		fmt.Printf("Uploaded %s\n", req.Filename)

	case protocol.CmdDelete:
		fmt.Printf("Deleted %s\n", req.Filename)

	case protocol.CmdGetFilesList:
		if len(res.Files) == 0 {
			fmt.Println("(empty directory)")
		} else {
			for _, f := range res.Files {
				fmt.Printf("%s\n", f)
			}
		}

	case protocol.CmdGetFile:
		target := filepath.Join(c.downloadDir, req.Filename)

		err := os.WriteFile(target, res.File, os.ModePerm)

		if err != nil {
			fmt.Printf("Error downloading file: %v\n", err)
		} else {
			fmt.Printf("Saved %s on %s\n", req.Filename, target)
		}
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading enviroment")
	}

	addr := os.Getenv("SERVER_ADDR")
	downloadDir := os.Getenv("DOWNLOAD_DIR")

	if addr == "" || downloadDir == "" {
		log.Fatal("Provide a valid enviroment")
	}

	err = os.MkdirAll(downloadDir, os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating download directory: %v", err)
	}

	client := NewClient(addr, downloadDir)
	client.Run()
}
