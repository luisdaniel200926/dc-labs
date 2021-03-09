/*
Clock Server is a concurrent TCP server that writes the time of a given timezone.
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func TimeIn(t time.Time, timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func handleConn(c net.Conn, timeZone string) {
	defer c.Close()
	t, err := TimeIn(time.Now(), timeZone)
	if err == nil {
		fmt.Println("Error")
	} else {
		fmt.Println(timeZone, "")
	}
	var locationAndTime = timeZone + "\t" + t.Format("15:04:05\n")
	_, er := io.WriteString(c, locationAndTime)
	if er != nil {
		fmt.Println("Error, client is not connected")
		return
	}

}

func main() {

	var port string
	flag.StringVar(&port, "port", "8080", "Eg: 9090")
	flag.Parse()
	port = "localhost:" + port
	var env = os.Getenv("TZ")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, env)
	}
}
