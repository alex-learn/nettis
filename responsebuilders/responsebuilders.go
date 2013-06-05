package responsebuilders

import( 
	"bytes"
	"encoding/hex"
	"log"
	"strings"
	)

type ResponseBuilderParams struct {
	Input []byte
}

type ResponseBuilder interface {
	GetResponse (params ResponseBuilderParams) ([]byte, error)
}

type EchoResponseBuilder struct {
}

func (r EchoResponseBuilder) GetResponse (params ResponseBuilderParams) ([]byte, error) {
	return params.Input, nil
}

type PrefixBasedResponseBuilder struct {
	StockResponses map[string]string
	DefaultResponse string
	RequestFilter string
	ResponseFilter string
}

func spaceyHexDecodeString(value string) ([]byte, error) {
	return hex.DecodeString(strings.NewReplacer(" ", "", "\n", "", "\r", "", "\t", "").Replace(value))
}

func applyFilter(filterName string, input string) (byts []byte, err error) {
	switch filterName {
	case "UNHEX":
		byts, err = spaceyHexDecodeString(input)
		if err != nil {
				return nil, err
		}
	case "HEX":
		byts = []byte(hex.EncodeToString([]byte(input)))
	default:
		byts = []byte(input)
	}
	return byts, err
}

func (r PrefixBasedResponseBuilder) GetResponse (params ResponseBuilderParams) ([]byte, error) {
	
	for k, v := range r.StockResponses {
		kBytes, err := applyFilter(r.RequestFilter, k)
		if err != nil {
				return nil, err
		}
		if bytes.HasPrefix(params.Input, kBytes) {
			o, err := applyFilter(r.ResponseFilter, v)
			return o, err
			log.Printf("Input:  %s -> Output: %s", params.Input, o)
		}
	}
	o, err := applyFilter(r.ResponseFilter, r.DefaultResponse)
	log.Printf("Input:  %s -> Output: %s", params.Input, o)
	return o, err
}