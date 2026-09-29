package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// classe Client (ou mantém sem OO, apenas impelementa hash pra enviar a senha no CONNECT e umas funções pra formatar prints)

// construtor

// run()

// hash

func main() {
	const addr = "127.0.0.1:9090"
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal("Error connecting: ", err)
	}

	defer conn.Close()
	fmt.Printf("Connected on %s\n", addr)

	inputReader := bufio.NewReader(os.Stdin)
	responseReader := bufio.NewReader(conn)

	for {
		fmt.Print(">> ")

		str, err := inputReader.ReadString('\n')

		if err != nil {
			log.Fatal("Error reading input: ", err)
		}

		fmt.Fprintf(conn, str)

		resp, err := responseReader.ReadString('\n')

		if err != nil {
			log.Fatal("Error reading response: ", err)
		}

		fmt.Print(resp)

		if strings.TrimSpace(string(str)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
