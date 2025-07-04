package controller

import (
	"log"
	"net/http"
	"test_go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub         *Hub
	authService service.AuthorService // или другие сервисы, если нужны
}

func NewWebSocketHandler(hub *Hub, authService service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := NewClient(conn, h.hub)
	h.hub.register <- client

	// Запускаем чтение и запись в отдельные goroutines
	go client.writePump()
	go client.readPump()
}
