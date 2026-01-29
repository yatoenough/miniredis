package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/yatoenough/miniredis/internal/handler"
	"github.com/yatoenough/miniredis/internal/resp"
	"github.com/yatoenough/miniredis/internal/writer"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Listening on port :6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		respr := resp.NewRESP(conn)
		value, err := respr.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		w := writer.NewWriter(conn)

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			w.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		result := handler(args)
		w.Write(result)
	}
}
