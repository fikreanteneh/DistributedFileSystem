package api

import (
	"dfs/internal/config"
	"dfs/internal/service"
	"fmt"
	"net"
	"net/http"
)

type ChunkServer struct {
	environment *config.Environment
	service     *service.ChunkService
}

func NewChunkServer(environment *config.Environment, service *service.ChunkService) *ChunkServer {
	return &ChunkServer{environment, service}
}

func (server *ChunkServer) run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/uploadChunk", server.UploadChunkHandler)
	mux.HandleFunc("/getChunk", server.GetChunkHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", server.environment.ChunkServerPort))
	if err != nil {
		panic(err)
	}
	http.Serve(listener, mux)
}

func (server *ChunkServer) UploadChunkHandler(w http.ResponseWriter, r *http.Request) {}
func (server *ChunkServer) GetChunkHandler(w http.ResponseWriter, r *http.Request)    {}
