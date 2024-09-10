package service

import (
	"dfs/internal/config"
	"dfs/internal/models"
	"dfs/internal/utils"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ChunkService struct {
	environment *config.Environment
}

func NewChunkService(env *config.Environment) *ChunkService {
	return &ChunkService{environment: env}
}

func (service *ChunkService) Upload() (any, error) {
	panic("not implemented")
}

func (service *ChunkService) Get() (any, error) {
	panic("not implemented")
}

func (service *ChunkService) HeartBeatEmit() error {
	url := fmt.Sprintf("ws://localhost:%s/heartbeat", service.environment.MasterServerPort)
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("error establishing WebSocket connection: %v", err)
		return err
	}
	defer wsConn.Close()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		heartbeatReq := models.HeartbeatNotifier{
			Url:  fmt.Sprintf(service.environment.ChunkServerURL),
			Port: service.environment.ChunkServerPort,
			Ip:   utils.GetOutboundIP(),
		}
		err := wsConn.WriteJSON(heartbeatReq)
		if err != nil {
			log.Printf("error sending heartbeat to master: %v", err)
			continue
		}
	}
	return nil
}

func (chunkService *ChunkService) SendSuccessfulUpload() (any, error) {
	panic("not implemented")
}
