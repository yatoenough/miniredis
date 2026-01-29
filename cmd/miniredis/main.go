package main

import (
	"fmt"
	"log"
	"net"

	"github.com/yatoenough/miniredis/internal/resp"
	"github.com/yatoenough/miniredis/internal/writer"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Listening on port :6379")

	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("%v", err)
	}

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

		writer := writer.NewWriter(conn)

		writer.Write(resp.Value{Typ: "string", Str: "OK"})
	}
}
