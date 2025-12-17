// Package util provides utility functions for the goku-framework.
package util

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*websocket.Conn]bool
	mu         sync.Mutex
	upgrader   websocket.Upgrader
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		upgrader:   websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections for development
			},
		},
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					// Assume client disconnected, schedule for removal
					go func(c *websocket.Conn) { h.unregister <- c }(conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	h.register <- conn

	// This is a one-way communication (server -> client),
	// but we need to keep the connection alive.
	// The client closing the connection will be handled by the write error in Run().
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	h.unregister <- conn
}

// BroadcastReload sends a "reload" message to all connected clients.
func (h *Hub) BroadcastReload() {
	log.Println("Change detected, broadcasting reload message to clients.")
	h.broadcast <- []byte("reload")
}

// WatchFiles starts watching files and triggers a reload on change.
func WatchFiles(hub *Hub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// We only care about write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					hub.BroadcastReload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	// Add directories/files to watch
	// For simplicity, we watch the entire project directory.
	// In a real application, you might want to be more specific.
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return watcher.Add(path)
	})

	if err != nil {
		log.Fatalf("Failed to add paths to watcher: %v", err)
	}

	// Block forever
	<-make(chan struct{})
}
