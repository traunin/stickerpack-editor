package main

import (
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/config"
	"github.com/Traunin/stickerpack-editor/apps/api/server"
)


func main() {
	config := config.Load()
	router := server.SetupRouter()
	
	addr := ":" + config.Port
	log.Printf("Server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}

	log.Println("Exiting...")
}
