package config

import (
	"github.com/laher/nettis/responsebuilders"
)

const MESSAGE_DEFAULT= "HELLO"

type Settings struct {
	Listen bool
	Http bool
	ResponseGenerator responsebuilders.ResponseBuilder
	//http/port
	Target string
	
	Initiate bool
	InitiateMessage string
	
	Delay int
	
	Verbose bool
	
	MaxReconnects int
	MaxMessages int
	
	Tls bool
	CertName string
	KeyName string
	TrustedCertName string
}

