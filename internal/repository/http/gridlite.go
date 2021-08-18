package http

import (
	"encoding/json"
	"log"

	// "os"
	// "time"

	"github.com/gorilla/websocket"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

type websocketGridLiteRepo struct {
	Conn *websocket.Conn
	// Channel chan
}

func NewGridLiteRepository(conn *websocket.Conn) repository.GridLite {

	log.Printf("executing go routine")
	return &websocketGridLiteRepo{
		Conn: conn,
	}
}
func (r *websocketGridLiteRepo) GetMessages(channel chan map[string]interface{}) {
	done := make(chan struct{})
	// interrupt := make(chan os.Signal, 1)
	go func() {
		defer close(done)
		for {
			log.Println("show")
			_, message, err := r.Conn.ReadMessage()
			var resp map[string]interface{}
			_ = json.Unmarshal(message, &resp)
			channel <- resp
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}
