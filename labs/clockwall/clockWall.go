package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var connections = os.Args[1:]
	var ports []string
	//fmt.Println(connections)
	for _, v := range connections {
		Splitted := strings.Split(v, "=")
		ports = append(ports, Splitted[1])
	}
	c := make(chan int)
	//fmt.Println(ports)
	for _, v := range ports {
		go printHour(v, c)
	}
	info := <-c
	fmt.Println(info)
	close(c)
}

func printHour(v string, c chan int) {
	conn, err := net.Dial("tcp", v)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, error2 := io.Copy(os.Stdout, conn)
	if error2 != nil {
		log.Fatal(error2)
	}
}
