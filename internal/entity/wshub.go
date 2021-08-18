package entity

import (
	"github.com/gorilla/websocket"
)

// Client is struct for handling websocket client connection
type Client struct {

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send         chan []byte
	StreamID     string
	GridLiteConn *websocket.Conn
}

// Message is struct for handling message received from websocket server
type Message struct {
	AnalyticID    string      `json:"analytic_id"`
	Image         string      `json:"image_jpeg,omitempty"`
	NodeNum       int64       `json:"node_num"`
	PipelineData  interface{} `json:"pipeline_data"`
	PrimaryText   string      `json:"primary_text"`
	SecondaryText string      `json:"secondary_text"`
	StreamAddress string      `json:"stream_address"`
	StreamID      string      `json:"stream_id"`
	Timestamp     int64       `json:"timestamp"`
}

// func (c *entity.Client) WritePump(s *ServiceImpl) {
// 	for {
// 		select {
// 		case message, ok := <-c.Send:

// 			// fmt.Println("message cummmm", message)
// 			var jsonData map[string]interface{}
// 			err := json.Unmarshal(message, &jsonData)
// 			if err != nil {
// 				log.Println("fail to parse", err)
// 			}
// 			fmt.Println(jsonData, "jsonData")

// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			pipelineData := jsonData["pipeline_data"].(map[string]interface{})
// 			if pipelineData["status"] == nil {
// 				w.Write(message)
// 			} else if pipelineData["status"].(string) != "UNKNOWN" {
// 				additionalData, errAdd := s.AddAdditionalData(jsonData)
// 				if errAdd != nil {
// 					w.Write(message)

// 				} else {
// 					byteData, _ := json.Marshal(additionalData)
// 					w.Write(byteData)
// 				}
// 			} else {
// 				w.Write(message)
// 			}

// 			// Add queued chat messages to the current websocket message.
// 			n := len(c.Send)
// 			for i := 0; i < n; i++ {
// 				w.Write(newline)
// 				w.Write(<-c.Send)
// 			}

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
