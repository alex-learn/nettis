package transport

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	RECV_BUF_LEN = 1024
)

func Connect(port string, initiate bool, delay int, verbose bool) {
	addr := port
	if strings.Index(port, ":") < 0 {
		addr = "127.0.0.1:" + port
	}
	log.Printf("Starting connection to %s", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("error connecting: %s", err.Error())
		os.Exit(1)
	}
	//same thread (prevent program exit)
	EchoService(conn, initiate, delay, verbose)
	log.Printf("Finished ")
}

func Listen(port string, initiate bool, delay int, verbose bool) error {
	addr := port
	if strings.Index(port, ":") < 0 {
		addr = "0.0.0.0:" + port
	}
	log.Printf("Starting server on %s", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error listening: %s", err.Error())
		return err
	}
	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accept: %s", err.Error())
			return err
		}
		go EchoService(conn, initiate, delay, verbose)
		i = i + 1
	}
	return nil
}
func write(conn net.Conn, buf []byte, delay int, verbose bool) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}
	n, err := conn.Write(buf)
	if err != nil {
		log.Printf("Error send: %s", err.Error())
		conn.Close()
	} else {
		if verbose {
			fmt.Printf("write: %s\n", buf)
		} else {
			fmt.Printf("w%d ", n)
		}
	}
}

func EchoService(conn net.Conn, initiate bool, delay int, verbose bool) {
	defer conn.Close()
	log.Printf("New connection with: %s", conn.RemoteAddr().String())
	buf := make([]byte, RECV_BUF_LEN)
	if initiate {
		go write(conn, buf, delay, verbose)
	}

	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("Error reading: %s", err.Error())
		conn.Close()
		return
	} else {
		if verbose {
			fmt.Printf("read: %s\n", buf[0:n])
		} else {
			fmt.Printf("r%d ", n)
		}
	}
	for err == nil {
		go write(conn, buf[0:n], delay, verbose)
		n, err = conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading: %s", err.Error())
			conn.Close()
			return
		} else {
			if verbose {
				fmt.Printf("read: %s\n", buf[0:n])
			} else {
				fmt.Printf("r%d ", n)
			}
		}
	}
	if err != nil {
		log.Printf("Read error: %s", err.Error())
		conn.Close()
	}
	log.Printf("EchoService finished")
}
