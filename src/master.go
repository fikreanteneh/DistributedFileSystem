package main

import (
	models "dfs/src/models"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


type MasterServer struct {
	Port            string
	ChunkSize       uint64
	ChunkserverUrls map[string]time.Time
	RPCPort         string
	Mutex           sync.Mutex
}

func NewMasterServer(port string, chunkSize uint64) *MasterServer {
	return &MasterServer{port, chunkSize, make(map[string]time.Time), "7000" ,sync.Mutex{}}
}

func (master *MasterServer) run() error {
	handler := enableCORS(http.DefaultServeMux)
	go master.registerUploadEndpoint()
	go master.registerGetEndpoint()
	// go master.registerReportChunkUploadSuccessEndpoint()
	// master.registerRegisterChunkserverEndpoint()
	go master.registerListenToHeartbeat()
	go master.registerViewEndpoint()
	go func() {
        for {
            time.Sleep(5 * time.Second) // Check every 5 seconds
            currentTime := time.Now()
            for url, lastHeartbeat := range master.ChunkserverUrls {
				if currentTime.Sub(lastHeartbeat) >= 5*time.Second {
					master.Mutex.Lock()
                    log.Println("No heartbeat received from URL:", url)
                    delete(master.ChunkserverUrls, url)
					master.Mutex.Unlock()
                }
            }
        }
    }()
	go master.RegisterRPCMethods()
	if err := http.ListenAndServe(":"+master.Port, handler); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (master *MasterServer) registerUploadEndpoint() {
	http.HandleFunc("/upload", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		log.Print("==========upload endpoint hit===================")
		var uploadReq models.UploadInitRequest
		if err := json.NewDecoder(r.Body).Decode(&uploadReq); err != nil {
			http.Error(w, "error parsing request body", 400)
			return
		}

		if len(master.ChunkserverUrls) == 0 {
			http.Error(w, "no chunkserver available", 500)
			return
		}

		fileIdentifier := getIdentifierFromFilename(uploadReq.FileName)
		numberOfChunks := calcNumberOfChunks(uploadReq.FileSize, master.ChunkSize)

		// urls := reflect.ValueOf(master.ChunkserverUrls).MapKeys()
		// randomChunkserver := urls[rand.Intn(len(urls))].String()
		chunkservers := make([]string, 0, len(master.ChunkserverUrls))
		for k := range master.ChunkserverUrls {
			chunkservers = append(chunkservers, k)
		}

		file := models.FileMetadata{
			Id: primitive.NewObjectID(),
			FileName: uploadReq.FileName, 
			FileSize: uploadReq.FileSize, 
			NumberOfChunks: numberOfChunks, 
			Replicas: make([][]string, numberOfChunks), 
			ClientId: "", 
			SharedUser: []primitive.ObjectID{},
			FileIdentifier: fileIdentifier,
		}
		log.Print("==========Before Create ======")
		file.CreateFile()
		log.Print("==========After Create ======")

		response := models.UploadInitResponse{
			Identifier:     fileIdentifier,
			ChunkSize:      master.ChunkSize,
			NumberOfChunks: numberOfChunks,
			Chunkservers:   chunkservers,
		}
		json.NewEncoder(w).Encode(response)
	}))
}

func isPortAlive(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}


func (master *MasterServer) isActiveChunk(url string) bool {
	master.Mutex.Lock()
    defer master.Mutex.Unlock()
    _, exists := master.ChunkserverUrls[url]
    return exists
}

func (master *MasterServer) registerGetEndpoint() {
	http.HandleFunc("/get", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fileIdentifier := r.URL.Query()["id"][0]
		// metadata, ok := master.Files[fileIdentifier]
		metadata, ok := models.GetFileById(fileIdentifier)

		if ok != nil{
			http.Error(w, fmt.Sprintf("file with identifier '%v' not found", fileIdentifier), http.StatusNotFound)
			return
		}

		log.Print("==========replica======", metadata.Replicas)
		locations := make([]string, metadata.NumberOfChunks)

		var i []string = metadata.Replicas[0]
		var available string = ""
		for _, value := range i {
			if master.isActiveChunk(value) {
				available = value
				break
			}

			// check if localhost:8000 is alive or not down
			// response, err := http.Get(value + "get")
			// if err == nil && response.StatusCode == http.StatusOK {
			// 	// If the request was successful, consider it available
			// 	available = value
			// 	break
			// }

			// response, err := http.Get(value)
			// if response.StatusCode == http.StatusOK && err == nil {
			// 	break
			// }

		}
		for i := range metadata.Replicas {
			locations[i] = available
		}

		response := models.GetResponse{
			FileName:  metadata.FileName,
			Locations: locations,
		}
		json.NewEncoder(w).Encode(response)
	}))
}

// func (master *MasterServer) registerReportChunkUploadSuccessEndpoint() {
// 	http.HandleFunc("/uploadSuccessful", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
// 		var uploadSuccessReq models.ChunkUploadSuccessRequest
// 		if err := json.NewDecoder(r.Body).Decode(&uploadSuccessReq); err != nil {
// 			http.Error(w, "error parsing request body", 400)
// 			return
// 		}

// 		parts := strings.Split(uploadSuccessReq.ChunkIdentifier, "_")
// 		fileIdentifier := parts[0]
// 		chunkIndex, err := strconv.Atoi(parts[1])
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("error parsing chunk index: %v", err), http.StatusBadRequest)
// 		}

// 		file, err := models.GetFileById(fileIdentifier)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("error fetching file metadata: %v", err), http.StatusBadRequest)
// 		}
// 		file.Replicas[chunkIndex] = append(file.Replicas[chunkIndex], uploadSuccessReq.Chunkserver)
// 		file.UpdateFile()

// 		// master.Files[fileIdentifier].Replicas[chunkIndex] = append(master.Files[fileIdentifier].Replicas[chunkIndex], uploadSuccessReq.Chunkserver)
// 	}))
// }




// Define your RPC server struct
type RPCServer struct {
	Master *MasterServer
}

// Define a method on RPCServer for handling the RPC call
func (rpc *RPCServer) ReportChunkUploadSuccess(req models.ChunkUploadSuccessRequest, reply *string) error {
	parts := strings.Split(req.ChunkIdentifier, "_")
	fileIdentifier := parts[0]
	chunkIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("error parsing chunk index: %v", err)
	}

	file, err := models.GetFileById(fileIdentifier)
	if err != nil {
		return fmt.Errorf("error fetching file metadata: %v", err)
	}
	file.Replicas[chunkIndex] = append(file.Replicas[chunkIndex], req.Chunkserver)
	file.UpdateFile()

	// You can set a response if necessary
	*reply = "Success"
	
	return nil
}

// This function registers the RPC method to handle incoming requests
func (master *MasterServer) RegisterRPCMethods() {
	rpcServer := new(RPCServer)
	rpcServer.Master = master
	err := rpc.Register(rpcServer)
	rpcListener, err := net.Listen("tcp", ":"+master.RPCPort)
	rpc.HandleHTTP()
	http.Serve(rpcListener, nil)
	if err != nil {
		log.Fatal("Register RPC Server error:", err)
	}
}










// func (master *MasterServer) registerRegisterChunkserverEndpoint() {
// 	http.HandleFunc("/chunkserver", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
// 		log.Print("==========register chunkserver endpoint hit===================")
// 		var registerReq models.RegisterChunkserverRequest
// 		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
// 			http.Error(w, "error parsing request body", 400)
// 		}

// 		master.ChunkserverUrls[registerReq.Url] = time.Now()

// 		response := models.RegisterChunkserverResponse{}
// 		json.NewEncoder(w).Encode(response)
// 	}))
// }


func (master *MasterServer) registerListenToHeartbeat() {
    http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println(err)
            return
        }
        defer conn.Close()

		 for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                log.Println("read failed:", err)
                return
            }
            var heartbeatMessage models.HeartbeatRequest
            if err := json.Unmarshal(message, &heartbeatMessage); err != nil {
                log.Println("error decoding heartbeat message:", err)
                return
            }

            // Extract the port from the heartbeat message
			master.Mutex.Lock()
			master.ChunkserverUrls[heartbeatMessage.Url] = time.Now()
			master.Mutex.Unlock()
            log.Println("=============== Received heartbeat from Location:", heartbeatMessage.Url, " ===============")
        }
        
    })
}

func (master *MasterServer) registerViewEndpoint() {
	http.HandleFunc("/view", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {

		files, err := models.GetAllFiles()
		if err != nil {
			http.Error(w, "error fetching file metadata", 400)
		}

		response := models.ViewResponse{
			Files: files,
		}
		json.NewEncoder(w).Encode(response)
	}))
}



func enableCORS(handler http.Handler) http.Handler { 
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
  // Allow requests from any origin 
  w.Header().Set("Access-Control-Allow-Origin", "*") 
 
  // Allow the GET, POST, and OPTIONS methods 
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") 
 
  // Allow the Content-Type header 
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type") 
 
  // If the request method is OPTIONS, return immediately with a 200 status code 
  if r.Method == "OPTIONS" { 
   w.WriteHeader(http.StatusOK) 
   return 
  } 
 
  // Call the next handler 
  handler.ServeHTTP(w, r) 
 }) 
}