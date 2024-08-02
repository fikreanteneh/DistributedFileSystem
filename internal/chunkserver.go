package main

import (
	"bytes"
	"dfs/src/models"

	// "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

type Chunkserver struct {
	Master string
	Port   string
	Dir    string
}

func NewChunkserver(master string, port string, dir string) *Chunkserver {
	return &Chunkserver{master, port, dir}
}

func (chunkserver *Chunkserver) run() error {
	// registerResponse, err := chunkserver.registerAtMaster()
	// if err != nil {
	// 	log.Fatalf("registering at master failed: %v", err)
	// }

	// log.Print(registerResponse)
	go chunkserver.registerUploadChunk()
	go chunkserver.registerGetChunkEndpoint()
	// chunkserver.registerChunkHome()
	go chunkserver.registerSendHeartbeatToMaster()

	if err := http.ListenAndServe(":"+chunkserver.Port, nil); err != nil {
		log.Fatal(err)
	}

	return nil
}


func (chunkserver *Chunkserver) registerUploadChunk() {
	http.HandleFunc("/uploadChunk", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("error parsing multipart form: %v", err), http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("chunk")
		if err != nil {
			http.Error(w, "error parsing chunk", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			http.Error(w, "error parsing chunk", http.StatusBadRequest)
			return
		}

		_, err = os.Create(filepath.Join(chunkserver.Dir, fileHeader.Filename))
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating file: %v", err), http.StatusInternalServerError)
			return
		}

		chunk := buf.Bytes()
		ioutil.WriteFile(filepath.Join(chunkserver.Dir, fileHeader.Filename), chunk, 0777)

		chunkserver.reportChunkUploadSuccessToMaster(fileHeader.Filename)
	}))
}

// func (chunkserver *Chunkserver) registerAtMaster() (*models.RegisterChunkserverResponse, error) {
// 	registerReq := models.RegisterChunkserverRequest{
// 		Url: fmt.Sprintf("%v:%v/", GetOutboundIP(), chunkserver.Port),
// 	}

// 	var registerResponse models.RegisterChunkserverResponse
// 	if err := postJson(chunkserver.Master+"chunkserver", &registerReq, &registerResponse); err != nil {
// 		return nil, err
// 	}

// 	return &registerResponse, nil
// }

// func (chunkserver *Chunkserver) reportChunkUploadSuccessToMaster(chunkIdentifier string) error {
// 	uploadSuccessfulReq := models.ChunkUploadSuccessRequest{chunkIdentifier, fmt.Sprintf("%v:%v/", GetOutboundIP(), chunkserver.Port)}
// 	client := &http.Client{Timeout: 10 * time.Second}
// 	requestBody, _ := json.Marshal(uploadSuccessfulReq)
// 	_, err := client.Post(chunkserver.Master+"uploadSuccessful", "application/json", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }




func (chunkserver *Chunkserver) reportChunkUploadSuccessToMaster(chunkIdentifier string) error {
    // Create a new RPC client
	log.Println("==============", fmt.Sprintf("%v:%v/", GetOutboundIP(), "7000"))
    client, err := rpc.DialHTTP("tcp", "localhost"+":7000")
	log.Println("==============", client)
    if err != nil {
        log.Fatal("Dialing:", err)
    }

	log.Println("==============Calling",)

    // Prepare the request
    uploadSuccessfulReq := models.ChunkUploadSuccessRequest{chunkIdentifier, fmt.Sprintf("%v:%v/", GetOutboundIP(), chunkserver.Port)}

    // Prepare a placeholder for the reply
    var reply string


    // Call the RPC method
    err = client.Call("RPCServer.ReportChunkUploadSuccess", uploadSuccessfulReq, &reply)
    if err != nil {
        log.Fatal("Upload error:", err)
    }
	log.Println("==============Called",)

    return nil
}



func (chunkserver *Chunkserver) registerGetChunkEndpoint() {
	http.HandleFunc("/get", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		chunkidentifier := r.URL.Query()["id"][0]

		f, err := os.Open(filepath.Join(chunkserver.Dir, chunkidentifier))
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't open chunk file: %v", err), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, fmt.Sprintf("error writing chunk data into stream: %v", err), http.StatusInternalServerError)
			return
		}
	}))
}

func (chunkserver *Chunkserver) registerSendHeartbeatToMaster() {
	    var err error
	url :=	"ws://localhost:8000/heartbeat"
    wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatalf("error establishing WebSocket connection: %v", err)
    }
    defer wsConn.Close()
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        heartbeatReq := models.HeartbeatRequest{
            Url: fmt.Sprintf("%v:%v/", GetOutboundIP(), chunkserver.Port),
			Port: chunkserver.Port,
			Ip: GetOutboundIP(),
        }
        err := wsConn.WriteJSON(heartbeatReq)
        if err != nil {
            log.Printf("error sending heartbeat to master: %v", err)
            continue
        }
    }
}






