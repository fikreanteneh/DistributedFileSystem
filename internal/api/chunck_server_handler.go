package api

import (
	"dfs/internal/config"
	chunkservice "dfs/internal/service/chunk_service"
	"fmt"
	"log"
	"net"
	"net/http"
)

type ChunkServer struct {
	environment *config.Environment
	service     *chunkservice.ChunkService
}

func NewChunkServer(environment *config.Environment, service *chunkservice.ChunkService) *ChunkServer {
	return &ChunkServer{environment, service}
}

func (server *ChunkServer) run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/uploadChunk", server.UploadChunkHandler)
	mux.HandleFunc("/getChunk", server.GetChunkHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", server.environment.ChunkServerPort))
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	http.Serve(listener, mux)
	return nil
}

func (server *ChunkServer) UploadChunkHandler(w http.ResponseWriter, r *http.Request) {}
func (server *ChunkServer) GetChunkHandler(w http.ResponseWriter, r *http.Request)    {}
