package main

import (
	"fmt"
	"github.com/laher/nettis/config"
	"github.com/laher/nettis/responsebuilders"
	"github.com/laher/nettis/transport"
	"github.com/laher/uggo"
	"log"
	"os"
)

const (
	CONFIG_FILE_DEFAULT = "nettis.json"
)
// VERSION is initialised by the linker during compilation if the appropriate flag is specified:
// e.g. go build -ldflags "-X main.VERSION 0.1.2-abcd" goxc.go
// thanks to minux for this advice
// So, goxc does this automatically during 'go build'
var (
	VERSION		= "0.1.2"
	flagSet         = uggo.NewFlagSetDefault("nettis", "[options] [host:]<port>", VERSION)
	certName        string
	keyName         string
	trustedcertName string
	configFile      string
	settings        config.Settings
)

func fileExists(path string) (bool, error) {
        _, err := os.Stat(path)
        if err == nil {
                return true, nil
        }
        if os.IsNotExist(err) {
                return false, nil
        }
        return false, err
}

func main() {
	call := os.Args
	log.SetPrefix("[nettis] ")
	//flagSet.Record([]string{"version"}, "bool")
	//flagSet.BoolVar(&isVersion, "version", false, "Show version")
	flagSet.AliasedStringVar(&settings.InitiateMessage, []string{"im", "message"}, config.MESSAGE_DEFAULT, "initiating message")
	flagSet.AliasedBoolVar(&settings.Verbose, []string{"v", "verbose"}, false, "verbose")
	flagSet.AliasedBoolVar(&settings.Listen, []string{"l", "listen"}, false, "listen")
	flagSet.AliasedBoolVar(&settings.Initiate, []string{"i", "initiate"}, false, "initiate conversation")
	flagSet.AliasedIntVar(&settings.Delay, []string{"d", "delay"}, 0, "delay (seconds) before responding")
	flagSet.AliasedIntVar(&settings.MaxReconnects, []string{"r", "cr"}, 0, "(client only) max reconnections after a disconnection")
	flagSet.AliasedBoolVar(&settings.Http, []string{"w", "http"}, false, "Use HTTP (only server implemented so far)")
	flagSet.AliasedBoolVar(&settings.Tls, []string{"s", "ssl", "tls"}, false, "Secure sockets (TLS/SSL)")
	flagSet.StringVar(&settings.CertName, "s-cert", "cert.pem", "Certificate to use for TLS")
	flagSet.StringVar(&settings.TrustedCertName, "s-trusted-cert", "", "Trusted certificate to accept TLS (nil means trust-all)")
	flagSet.StringVar(&settings.KeyName, "s-key", "key.pem", "Key to use for TLS")
	flagSet.AliasedStringVar(&configFile, []string{"c", "config"}, CONFIG_FILE_DEFAULT, "Config file name")

	//TODO: cert config
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		os.Exit(1)
	}
	if flagSet.ProcessHelpOrVersion() {
		return
	}

	remainder := flagSet.Args()
	if len(remainder) < 1 {
		fmt.Fprintf(os.Stderr, "Error: not enough args\n")
		flagSet.Usage()
		os.Exit(1)
	}
	settings.Target = remainder[0]

	checkConfig := true

	//if no default config, use EchoResponseBuilder.
	//(if another config file specified, always try to load it (nettis will stop on error)
	if configFile == CONFIG_FILE_DEFAULT  {
		if exists, err := fileExists(configFile); exists {
			//OK
		} else {
			if err != nil {
				log.Fatalf("File check error: %s", err)
			}
			checkConfig = false
		}
	}
	if checkConfig {
		loadConfig(configFile, &settings)
	} else {
		log.Printf("'Echo' response generator")
		//default to EchoResponseBuilder
		settings.ResponseGenerator = responsebuilders.EchoResponseBuilder{}
	}
	
	if settings.Http {
		if settings.Listen {
			if settings.Tls {
				transport.ListenHttps(settings)
			} else {
				transport.ListenHttp(settings)
			}
		} else {
			if settings.Tls {
				log.Fatal("HTTPS client not implemented yet")
			} else {
				log.Printf("HTTP client")
				transport.ConnectHttp(settings)
			}
		}
	} else if settings.Tls {
		if settings.Listen {
			transport.ListenTls(settings)
		} else {
			transport.ConnectTls(settings)
		}
	} else {
		if settings.Listen {
			transport.Listen(settings)
		} else {
			transport.Connect(settings)
		}
	}
}

func loadConfig(configFile string, settings *config.Settings) {
	log.Printf("use config file %s", configFile)
	responderConfig, err := config.LoadJsonFile(configFile, true)
	if typ, keyExists := responderConfig["ResponseBuilder"]; keyExists {
		log.Printf("Type: %s", typ)
		switch typ {
		case "Prefix":
			log.Printf("Prefix based response generator")
			rmap, err := config.ToMapStringString(responderConfig["ResponseMap"], "ResponseMap")
			if err != nil {
				log.Printf("Response map error '%v'", err)
			}
			responseDefault, err := config.ToString(responderConfig["ResponseDefault"], "ResponseDefault")
			if err != nil {
				log.Fatalf("Default Response error '%v'", err)
			}
			requestFilter, err := config.ToString(responderConfig["RequestFilter"], "RequestFilter")
			if err != nil {
				log.Printf("RequestFilter error '%v'", err)
			}
			responseFilter, err := config.ToString(responderConfig["ResponseFilter"], "ResponseFilter")
			if err != nil {
				log.Printf("ResponseFiltermap error '%v'", err)
			}
			settings.ResponseGenerator = responsebuilders.PrefixBasedResponseBuilder{rmap, 
				responseDefault, 
				requestFilter,
				responseFilter}
		case "":
			//default to EchoResponseBuilder
			log.Printf("'Echo' response builder")
			settings.ResponseGenerator = responsebuilders.EchoResponseBuilder{}
		default:
			log.Fatalf("Unrecognised ResponseBuilder Type '%v'", typ)
		}
	} else {
		//default to EchoResponseBuilder
		log.Printf("'Echo' response generator")
		settings.ResponseGenerator = responsebuilders.EchoResponseBuilder{}
	}
	if err != nil {
		log.Fatalf("Error parsing responder config - %v", err)
	}
}
