package server

import (
	"log"
	"net"
)

// classe Server
type Server struct {
	addr     string
	listener net.Listener
}

// "construtor" do Server
func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// "run()" do Server (análogo em java); listen no addr e despachar goroutines para cada conexão nova
func (s *Server) Run() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)

	if err != nil {
		log.Fatal("Erro durante Listen: %v", err)
	}

	defer s.listener.Close()

	// log.Printf("Servidor executando em %s...\n", s.addr)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Printf("Erro durante Accept: %v\n", err)
		}

		session := NewSession(conn)
		go session.handleConnection()
	}
}

// Métodos para os comandos...

// hash (ver se tem SHA512 no go); o servidor fica com a senha? ou só o hash da senha?

// definir users (duda e carol)

func main() {
	const addr = "127.0.0.1:9090"

	// const addr = "127.0.0.1:9090"

	// listener, err := net.Listen("tcp", addr)

	// if err != nil {
	// 	log.Fatal("Error listeing: ", err)
	// }

	// defer listener.Close()
	// // log.Println("Listening on ", addr)

	// for {
	// 	conn, err := listener.Accept()

	// 	if err != nil {
	// 		log.Println("Error accepting conn: ", err)
	// 		continue
	// 	}

	// 	go handleConnection(conn)
	// }
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)

// 	for {
// 		message, err := reader.ReadString('\n')

// 		if err != nil {
// 			log.Printf("Read error: %v", err)

// 			return
// 		}

// 		fmt.Println(message)

// 		ackMsg := strings.ToUpper(strings.TrimSpace(message))
// 		response := fmt.Sprintf("ACK: %s\n", ackMsg)
// 		_, err = conn.Write([]byte(response))

// 		if err != nil {
// 			log.Printf("Server write error: %v", err)
// 		}
// 	}
// }
