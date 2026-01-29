package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/yatoenough/miniredis/internal/aof"
	"github.com/yatoenough/miniredis/internal/handler"
	"github.com/yatoenough/miniredis/internal/resp"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("miniredis listening on port :6379")

	aof := initAof()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		go handler.HandleConn(conn, aof)
	}
}

func initAof() *aof.AOF {
	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer aof.Close()

	aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	return aof
}
