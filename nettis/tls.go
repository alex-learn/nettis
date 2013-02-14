package nettis

import (
   "crypto/rand"
   "crypto/tls"
   "log"
   "net"
   "os"
   "crypto/rsa"
   "crypto/x509"
   "crypto/x509/pkix"
   "encoding/pem"
   "math/big"
   "time"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
   _, err := os.Stat(path)
   if err == nil { return true, nil }
   if os.IsNotExist(err) { return false, err }
   return false, err
}

func ListenTls(port string, initiate bool, delay int, certname string, keyname string, verbose bool) {
   c,err:=exists(certname)
   k,err:=exists(keyname)
   if !c {
      log.Printf("Cert file doesnt exist! %s", err)
      GenKeyCert("127.0.0.1", certname, keyname)
   } else if !k {
      log.Printf("Key file doesnt exist! %s", err)
      GenKeyCert("127.0.0.1", certname, keyname)
   }

   cert, err := tls.LoadX509KeyPair(certname, keyname)
   if err != nil {
      log.Fatalf("TLS: loadkeys: %s", err)
   }
   config := tls.Config{Certificates: []tls.Certificate{cert}} //, ClientAuth: tls.RequireAnyClientCert}
   config.Rand = rand.Reader
   service := "0.0.0.0:"+port
   listener, err := tls.Listen("tcp", service, &config)
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
      TlsService(conn, initiate, delay, verbose)
      i = i + 1
    }
}

func TlsService(conn net.Conn, initiate bool, delay int, verbose bool) {
    tlscon, ok := conn.(*tls.Conn)
    if ok {
        log.Print("server: conn: type assert to TLS succeedded")
        err := tlscon.Handshake()
        if err != nil {
            log.Printf("TLS ERROR: handshake failed: %s", err)
	    conn.Close()
        } else {
            log.Printf("TLS: conn: Handshake completed")
            go EchoService(conn, initiate, delay, verbose)
        }
    }
}




//Based on Go's source (http://golang.org/src/pkg/crypto/tls/generate_cert.go)
//Original copyright and licence ref below

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

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
