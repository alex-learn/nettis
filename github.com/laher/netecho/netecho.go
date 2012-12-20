package main

import (
   "flag"
   "os"
   "github.com/laher/netecho/netecho"
   "log"
   "fmt"
)
const NETECHO_VERSION="0.0.1"
var (
   flagSet    = flag.NewFlagSet("netecho", flag.PanicOnError)
   verbose bool
   listen bool
   initiate bool
   delay int
   http bool
   tls bool
   is_help bool
   is_version bool
   certname string
   keyname string
   trusted_certname string
)

func help_text() {
   fmt.Fprint(os.Stderr,"netecho [options] [host:]<port>\n")
   fmt.Fprintf(os.Stderr," Version %s. Options:\n", NETECHO_VERSION)
   flagSet.PrintDefaults()
}

func version_text() {
   fmt.Fprintf(os.Stderr," netecho version %s\n", NETECHO_VERSION)
}

func main() {
   call := os.Args
   log.SetPrefix("[netecho] ")
   flagSet.BoolVar(&verbose, "v", false, "verbose")
   flagSet.BoolVar(&listen, "l", false, "listen")
   flagSet.BoolVar(&initiate, "i", false, "initiate conversation")
   flagSet.IntVar(&delay, "d", 0, "delay (seconds) before echoing")
   flagSet.BoolVar(&http, "http", false, "initiate conversation")
   flagSet.BoolVar(&tls, "s", false, "Secure sockets (TLS/SSL)")
   flagSet.StringVar(&certname, "s-cert", "cert.pem", "Certificate to use for TLS")
   flagSet.StringVar(&trusted_certname, "s-trusted-cert", "", "Trusted certificate to accept TLS (nil means trust-all)")
   flagSet.StringVar(&keyname, "s-key", "key.pem", "Key to use for TLS")
   flagSet.BoolVar(&is_version, "version", false, "Show version")
   flagSet.BoolVar(&is_help, "h", false, "Show this help")
   //TODO: cert config
   e := flagSet.Parse(call[1:])
   if e != nil {
      os.Exit(1)
   }
   if is_help {
      help_text()
      return
   } else if is_version {
      version_text()
      return
   }
   remainder := flagSet.Args()
   if(len(remainder) < 1 ) {
      help_text()
      os.Exit(1)
   }
   port := remainder[0]
   if http {
      if listen {
         if tls {
            netecho.ListenHttps(port, delay)
         } else {
            netecho.ListenHttp(port, delay)
         }
      } else {
         log.Printf("HTTP client unimplemented")
         os.Exit(1)
      }
   } else if tls {
      if listen {
         netecho.ListenTls(port, initiate, delay,certname,keyname)
      } else {
         log.Printf("TLS client unimplemented")
         os.Exit(1)
      }
   } else {
      if listen {
         netecho.Listen(port, initiate, delay)
      } else {
         netecho.Connect(port, initiate, delay)
      }
   }
}
