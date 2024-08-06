package api

import (
	"dfs/internal/config"
	models "dfs/internal/models"
	masterservice "dfs/internal/service/master_service"
	"dfs/internal/utils"
	"fmt"
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
	environment *config.Environment
	service     *masterservice.MasterService
}

type MasterServer struct {
	environment *config.Environment
	upgrader    *websocket.Upgrader
	service     *masterservice.MasterService
	rpcListener *RPCListener
}

func NewMasterServer(environment *config.Environment, service *masterservice.MasterService) *MasterServer {
	return &MasterServer{
		environment: environment,
		service:     service,
		upgrader:    &upgrader,
		rpcListener: &RPCListener{environment: environment, service: service},
	}
}

func (master *MasterServer) run() error {
	mux := http.NewServeMux()
	rpc.Register(master.rpcListener)
	mux.HandleFunc("/upload", master.UploadHandler)
	mux.HandleFunc("/view", master.ViewHandler)
	mux.HandleFunc("/get", master.GetHandler)
	mux.HandleFunc("/heartbeat", master.HeratBeatHandler)
	mux.Handle("/rpc", rpc.DefaultServer)

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", master.environment.MasterServerPort))
	if err != nil {
		panic(err)
	}
	log.Println("Server listening on port " + master.environment.MasterServerPort)
	http.Serve(listener, mux)
	return nil
}

func (master *MasterServer) UploadHandler(w http.ResponseWriter, r *http.Request) {
	var uploadInit, err = utils.ReadJson[models.UploadInitRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := master.service.Upload(uploadInit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, result)
}

func (master *MasterServer) ViewHandler(w http.ResponseWriter, r *http.Request) {
	result, err := master.service.View()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, result)
}

func (master *MasterServer) GetHandler(w http.ResponseWriter, r *http.Request) {
	var fileId = r.URL.Query().Get("file")
	if fileId == "" {
		http.Error(w, "file id is required", http.StatusBadRequest)
		return
	}
	result, err := master.service.Get(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, result)
}

func (master *MasterServer) HeratBeatHandler(w http.ResponseWriter, r *http.Request) {

}

func (rpc *RPCListener) UploadSuccessful(args *models.ChunkUploadSuccessRequest, reply *models.GetFileResponse) error {
	panic("not implemented")
}
