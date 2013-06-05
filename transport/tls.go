package transport

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/laher/nettis/config"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}


func ConnectTls(settings config.Settings) {
	c, err := exists(settings.CertName)
	if !c {
		log.Printf("Cert file doesnt exist! %s", err)
		GenKeyCert("127.0.0.1", settings.CertName, settings.KeyName)
	}  
	//TODO client certs
	
	addr := GetAddress(settings, "127.0.0.1")
	log.Printf("Starting TLS connection to %s", addr)
	i := -1
	//only one at once. -1 represents 'forever'
	conf := tls.Config{
		//Certificates: []tls.Certificate{cert},
		CipherSuites: []uint16{
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
		InsecureSkipVerify : true}
	conf.Rand = rand.Reader
	for i < settings.MaxReconnects || settings.MaxReconnects < 0 {
		conn, err := tls.Dial("tcp", addr, &conf)
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


//func ListenTls(port string, initiate bool, delay int, certname string, keyname string, verbose bool) {
func ListenTls(settings config.Settings) {
	c, err := exists(settings.CertName)
	k, err := exists(settings.KeyName)
	if !c {
		log.Printf("Cert file doesnt exist! %s", err)
		GenKeyCert("127.0.0.1", settings.CertName, settings.KeyName)
	} else if !k {
		log.Printf("Key file doesnt exist! %s", err)
		GenKeyCert("127.0.0.1", settings.CertName, settings.KeyName)
	}

	cert, err := tls.LoadX509KeyPair(settings.CertName, settings.KeyName)
	if err != nil {
		log.Fatalf("TLS: loadkeys: %s", err)
	}
	conf := tls.Config{Certificates: []tls.Certificate{cert}} //, ClientAuth: tls.RequireAnyClientCert}
	conf.Rand = rand.Reader
	service := GetAddress(settings, "0.0.0.0")
	listener, err := tls.Listen("tcp", service, &conf)
	if err != nil {
		log.Fatalf("TLS: listen error: %s", err)
	}
	log.Printf("TLS: listening on %s", service)
	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TLS: accept: %s", err)
			break
		}
		log.Printf("TLS: accepted from %s", conn.RemoteAddr())
		TlsService(conn, settings)
		i = i + 1
	}
}

func TlsService(conn net.Conn, settings config.Settings) {
	tlscon, ok := conn.(*tls.Conn)
	if ok {
		log.Print("server: conn: type assert to TLS succeedded")
		err := tlscon.Handshake()
		if err != nil {
			log.Printf("TLS ERROR: handshake failed: %s", err)
			conn.Close()
		} else {
			log.Printf("TLS: conn: Handshake completed")
			ch := make(chan int, 100)
			go EchoService(conn, settings, ch)
		}
	}
}

//Based on Go's source (http://golang.org/src/pkg/crypto/tls/generate_cert.go)
//Original copyright and licence ref below

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

func GenKeyCert(host string, certname string, keyname string) error {

	log.Printf("Attempting to generate key %s and cert %s", certname, keyname)
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
		return err
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   host,
			Organization: []string{"Acme Co"},
		},
		NotBefore: now.Add(-5 * time.Minute).UTC(),
		NotAfter:  now.AddDate(1, 0, 0).UTC(), // valid for 1 year.

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Printf("Failed to create certificate: %s", err)
		return err
	}

	certOut, err := os.Create(certname)
	if err != nil {
		log.Printf("failed to open %s for writing: %s", certname, err)
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("written %s", certname)

	keyOut, err := os.OpenFile(keyname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("failed to open %s for writing: %s", keyname, err)
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Printf("written %s", keyname)
	return nil
}
