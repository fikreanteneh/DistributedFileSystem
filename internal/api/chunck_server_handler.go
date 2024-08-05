package api

import (
	"log"
	"net"
	"net/http"
)

type ChunkServer struct {
	Master string
	Port   string
	Dir    string
}

func NewChunkServer(master string, port string, dir string) *ChunkServer {
	return &ChunkServer{master, port, dir}
}

func (chunkserver *ChunkServer) run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/uploadChunk", chunkserver.UploadChunkHandler)
	mux.HandleFunc("/getChunk", chunkserver.GetChunkHandler)

	listener, err := net.Listen("tcp", ":"+chunkserver.Port)
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	http.Serve(listener, mux)
	return nil
}

func (chunkserver *ChunkServer) UploadChunkHandler(w http.ResponseWriter, r *http.Request) {}
func (chunkserver *ChunkServer) GetChunkHandler(w http.ResponseWriter, r *http.Request)    {}
