package main

import (
	"flag"
	"fmt"
	"github.com/laher/nettis/transport"
	"log"
	"os"
)

// VERSION is initialised by the linker during compilation if the appropriate flag is specified:
// e.g. go build -ldflags "-X main.VERSION 0.1.2-abcd" goxc.go
// thanks to minux for this advice
// So, goxc does this automatically during 'go build'
var VERSION string

var (
	flagSet         = flag.NewFlagSet("nettis", flag.PanicOnError)
	verbose         bool
	listen          bool
	initiate        bool
	delay           int
	http            bool
	tls             bool
	isHelp          bool
	isVersion       bool
	certName        string
	keyName         string
	trustedcertName string
)

func printHelp() {
	fmt.Fprint(os.Stderr, "nettis [options] [host:]<port>\n")
	fmt.Fprintf(os.Stderr, " Version '%s'. Options:\n", VERSION)
	flagSet.PrintDefaults()
}

func printVersion() {
	fmt.Fprintf(os.Stderr, " nettis version: %s\n", VERSION)
}

func main() {
	call := os.Args
	log.SetPrefix("[nettis] ")
	flagSet.BoolVar(&verbose, "v", false, "verbose")
	flagSet.BoolVar(&listen, "l", false, "listen")
	flagSet.BoolVar(&initiate, "i", false, "initiate conversation")
	flagSet.IntVar(&delay, "d", 0, "delay (seconds) before echoing")
	flagSet.BoolVar(&http, "http", false, "initiate conversation")
	flagSet.BoolVar(&tls, "s", false, "Secure sockets (TLS/SSL)")
	flagSet.StringVar(&certName, "s-cert", "cert.pem", "Certificate to use for TLS")
	flagSet.StringVar(&trustedcertName, "s-trusted-cert", "", "Trusted certificate to accept TLS (nil means trust-all)")
	flagSet.StringVar(&keyName, "s-key", "key.pem", "Key to use for TLS")
	flagSet.BoolVar(&isVersion, "version", false, "Show version")
	flagSet.BoolVar(&isHelp, "h", false, "Show this help")
	//TODO: cert config
	e := flagSet.Parse(call[1:])
	if e != nil {
		os.Exit(1)
	}
	if isHelp {
		printHelp()
		return
	} else if isVersion {
		printVersion()
		return
	}
	remainder := flagSet.Args()
	if len(remainder) < 1 {
		printHelp()
		os.Exit(1)
	}
	port := remainder[0]
	if http {
		if listen {
			if tls {
				transport.ListenHttps(port, delay, verbose)
			} else {
				transport.ListenHttp(port, delay, verbose)
			}
		} else {
			log.Printf("HTTP client unimplemented")
			os.Exit(1)
		}
	} else if tls {
		if listen {
			transport.ListenTls(port, initiate, delay, certName, keyName, verbose)
		} else {
			log.Printf("TLS client unimplemented")
			os.Exit(1)
		}
	} else {
		if listen {
			transport.Listen(port, initiate, delay, verbose)
		} else {
			transport.Connect(port, initiate, delay, verbose)
		}
	}
}
