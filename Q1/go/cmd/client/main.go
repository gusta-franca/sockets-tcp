package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"sd/sockets_tcp/internal/auth"
)

// classe Client
type Client struct {
	addr           string
	conn           net.Conn
	inputReader    *bufio.Reader
	responseReader *bufio.Reader
}

// construtor
func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.addr)

	if err != nil {
		return err
	}

	c.inputReader = bufio.NewReader(os.Stdin)
	c.responseReader = bufio.NewReader(c.conn)

	fmt.Printf("Connected on %s\n", c.addr)
	return nil
}

// run()
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
			fmt.Fprintf(c.conn, "EXIT\n")
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

			cmd = fmt.Sprintf("CONNECT %s,%s\n", usr, auth.HashSHA512(pwd))

		} else if command == "CHDIR" {
			if len(cmdTokens) < 2 {
				fmt.Println("Usage: CHDIR directory")
				continue
			}

			cmd = fmt.Sprintf("%s %s\n", command, strings.TrimSpace(cmdTokens[1]))
		} else {
			cmd = fmt.Sprintf("%s\n", command)
		}

		fmt.Fprintf(c.conn, "%s", cmd)

		response, err := c.responseReader.ReadString('\n')

		if err != nil {
			log.Printf("Error reading response: %s", err)
			log.Printf("Disconnected from %s\n", c.addr)
			break
		}

		fmt.Print(response)

		if command == "GETFILES" || command == "GETDIRS" {
			count := strings.TrimSpace(response) // file/dir count

			if count != "ERROR" {
				itemCount, err := strconv.Atoi(count)

				if err == nil {
					for range itemCount {
						line, err := c.responseReader.ReadString('\n')

						if err != nil {
							break
						}

						fmt.Print(line)
					}
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
