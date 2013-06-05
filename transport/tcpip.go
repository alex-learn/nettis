package transport

import (
	"io"
	"fmt"
	"github.com/laher/nettis/config"
	"github.com/laher/nettis/responsebuilders"
	"log"
	"net"
	"strings"
	"time"
)

const (
	RECV_BUF_LEN = 1024
)

func Connect(settings config.Settings) {
	addr := GetAddress(settings, "127.0.0.1")
	log.Printf("Starting connection to %s", addr)
	i := -1
	//only one at once. -1 represents 'forever'
	for i < settings.MaxReconnects || settings.MaxReconnects < 0 {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Printf("error connecting: %s", err.Error())
		} else {
			if settings.Verbose {
				log.Printf("Connected")
			}
			//same thread (prevent program exit)
			ch := make(chan int, 100)
			EchoService(conn, settings, ch)
			foo, ok := <- ch
			log.Printf("response ok: %b, code: %d", ok, foo)
		}
		i++
	}
	log.Printf("Finished after %d connections", i+1)
}

func GetAddress(settings config.Settings, defaultHost string) string {
	addr := settings.Target
	if strings.Index(settings.Target, ":") < 0 {
		addr = defaultHost + ":" + settings.Target
	}
	return addr
}

func Listen(settings config.Settings) error {
	addr := GetAddress(settings, "0.0.0.0")
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
		ch := make(chan int, 100)
		go EchoService(conn, settings, ch)
		i = i + 1
	}
	return nil
}

func write(conn net.Conn, buf []byte, settings config.Settings, ch chan int) {
	if settings.Delay > 0 {
		time.Sleep(time.Duration(settings.Delay) * time.Second)
	}
	if settings.Verbose {
		log.Printf("Writing buffer")
	}
	n, err := conn.Write(buf)
	if err != nil {
		log.Printf("Error send: %s", err.Error())
		conn.Close()
		ch <- -2
	} else {
		if settings.Verbose {
			log.Printf("write: %s (%d)\n", buf, n)
		} else {
			//fmt.Printf("w%d ", n)
			fmt.Printf("w")
		}
	}
}

func EchoService(conn net.Conn, settings config.Settings, ch chan int) {
	defer conn.Close()
	log.Printf("New connection with: %s", conn.RemoteAddr().String())
	buf := make([]byte, RECV_BUF_LEN)
	if settings.Initiate {
		buf = []byte(settings.InitiateMessage)
		go write(conn, buf, settings, ch)
	}

	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("Error reading: %s", err.Error())
		conn.Close()
		ch <- -1
		return
	} else {
		if settings.Verbose {
			log.Printf("read: %s\n", buf[0:n])
		} else {
			//fmt.Printf("r%d ", n)
			fmt.Printf("r")
		}
	}
	for err == nil {
		resp, err := settings.ResponseGenerator.GetResponse(responsebuilders.ResponseBuilderParams{buf[0:n]})
		if err != nil {
			log.Printf("Error generating response: %s", err.Error())
			conn.Close()
			ch <- -1
			return
		}
		go write(conn, resp, settings, ch)
		if settings.Verbose {
			log.Printf("Reading buffer")
		}
		n, err = conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading: %s", err.Error())
			conn.Close()
			ch <- -1
			return
		} else {
			if settings.Verbose {
				log.Printf("read: %s\n", buf[0:n])
			} else {
				//fmt.Printf("r%d ", n)
				fmt.Printf("r")
			}
		}
	}
	if err != nil {
		log.Printf("Read error: %s", err.Error())
		conn.Close()
		ch <- -1
		return
	}
	log.Printf("EchoService finished")
	ch <- 0
}
