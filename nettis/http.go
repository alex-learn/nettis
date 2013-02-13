package nettis

import (
   "log"
   "net/http"
   "time"
   "strings"
)


func handler(w http.ResponseWriter, req *http.Request) {
   w.Header().Set("Content-Type", "text/plain")
   w.Write([]byte("Default ECHO Response.\n"))
}


// initiate not suppported
func ListenHttp(port string, delay int) {
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      time.Sleep(time.Duration(delay) * time.Second)
      //TODO echo type and content
      w.Header().Set("Content-Type", "text/plain")
      w.Write([]byte("Default ECHO Response.\n"))
   })
   addr := port
   if strings.Index(port, ":")<0 {
      addr= "0.0.0.0:"+port
   }
   log.Printf("About to listen on http://"+addr+"/")
   err := http.ListenAndServe(addr, nil)
   if err != nil {
      log.Fatal(err)
   }
}

func ListenHttps(port string, delay int) {
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      time.Sleep(time.Duration(delay) * time.Second)
      //TODO echo type and content
      w.Header().Set("Content-Type", "text/plain")
      w.Write([]byte("Default ECHO Response.\n"))
   })
   addr := port
   if strings.Index(port, ":")<0 {
      addr= "0.0.0.0:"+port
   }
   log.Printf("About to listen on https://"+addr+"/")
   err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
   if err != nil {
      log.Fatal(err)
   }
}
