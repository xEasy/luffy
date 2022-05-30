package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client starting")
	conn, err := net.Dial("tcp4", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("client dial fail err: ", err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello luffy"))
		if err != nil {
			fmt.Println("writ err : ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf err: ", err)
		}
		fmt.Printf("Server call back: %s, cnt: %d \n", buf, cnt)
		time.Sleep(time.Second)
	}
}
