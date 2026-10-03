/*
   Data de criação: 01/10/2026
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
	"strconv"
	"strings"

	"sockets_tcp/Q1/internal/auth"
	"sockets_tcp/Q1/internal/protocol"

	"github.com/joho/godotenv"
)

type Client struct {
	addr        string
	conn        net.Conn
	proto       *protocol.Utf8Protocol
	inputReader *bufio.Reader
}

// Client constructor
func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
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
	c.proto = protocol.NewUtf8Protocol(c.conn)

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

		cmdTokens := strings.SplitN(strings.TrimSpace(str), " ", 2)
		command := strings.ToUpper(cmdTokens[0])

		if command == "" {
			continue
		}

		if command == "EXIT" {
			_ = c.proto.WriteString("EXIT")
			fmt.Printf("Disconnected from %s\n", c.addr)
			return
		}

		var cmd string

		if command == "CONNECT" {
			// "few arguments"
			if len(cmdTokens) < 2 {
				fmt.Println("Usage: CONNECT user,password")
				continue
			}

			argsTokens := strings.SplitN(cmdTokens[1], ",", 2)

			//"too many arguments"
			if len(argsTokens) != 2 {
				fmt.Println("Usage: CONNECT user,password")
				continue
			}

			usr := strings.TrimSpace(argsTokens[0])
			pwd := strings.TrimSpace(argsTokens[1])

			cmd = fmt.Sprintf("CONNECT %s,%s", usr, auth.HashSHA512(pwd))
		} else if command == "CHDIR" {
			if len(cmdTokens) < 2 {
				fmt.Println("Usage: CHDIR directory")
				continue
			}

			cmd = fmt.Sprintf("%s %s", command, strings.TrimSpace(cmdTokens[1]))
		} else {
			cmd = command
		}

		err = c.proto.WriteString(cmd)

		if err != nil {
			log.Printf("Error sending command: %s", err)
			break
		}

		response, err := c.proto.ReadString()

		if err != nil {
			log.Printf("Error reading response: %s", err)
			log.Printf("Disconnected from %s\n", c.addr)
			break
		}

		if response != "" {
			fmt.Println(response)
		}

		if (command == "GETFILES" || command == "GETDIRS") && response != "ERROR" {
			count := strings.TrimSpace(response)
			itemCount, err := strconv.Atoi(count)

			if err == nil {
				for range itemCount {
					item, err := c.proto.ReadString()

					if err != nil {
						log.Printf("Error reading item: %s", err)
						break
					}
					fmt.Println(item)
				}
			}
		}
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading enviroment")
	}

	addr := os.Getenv("SERVER_ADDR")

	if addr == "" {
		log.Fatal("Provide a valid enviroment")
	}

	client := NewClient(addr)
	client.Run()
}
