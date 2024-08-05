package api

import (
	models "dfs/internal/models"
	masterservice "dfs/internal/service/master_service"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type RPCListener struct {
	service *masterservice.MasterService
}

type MasterServer struct {
	Port     string
	upgrader websocket.Upgrader
	service  *masterservice.MasterService
}

func NewMasterServer(port string, rpcPort string, upgrader websocket.Upgrader, service *masterservice.MasterService) *MasterServer {
	return &MasterServer{port, upgrader, service}
}

func (master *MasterServer) run() error {
	mux := http.NewServeMux()
	rpc.Register(RPCListener{service: master.service})

	mux.HandleFunc("/upload", master.UploadHandler)
	mux.HandleFunc("/view", master.ViewHandler)
	mux.HandleFunc("/get", master.GetHandler)
	mux.HandleFunc("/heartbeat", master.HeratBeatHandler)
	mux.Handle("/rpc", rpc.DefaultServer)

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+master.Port)
	if err != nil {
		log.Fatal("Listener error: ", err)
	}

	log.Println("Server listening on port " + master.Port)
	http.Serve(listener, mux)

	return nil
}

func (master *MasterServer) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// var file models.FileMetadata
}

func (master *MasterServer) ViewHandler(w http.ResponseWriter, r *http.Request) {
}

func (master *MasterServer) GetHandler(w http.ResponseWriter, r *http.Request) {
}

func (master *MasterServer) HeratBeatHandler(w http.ResponseWriter, r *http.Request) {
	// var heartbeatMessage models.HeartbeatRequest
	// err := json.NewDecoder(r.Body).Decode(&heartbeatMessage)
	panic("not implemented")
}

func (rpc *RPCListener) UploadSuccessful(args *models.ChunkUploadSuccessRequest, reply *models.GetResponse) error {
	panic("not implemented")
}
