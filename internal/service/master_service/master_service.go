package masterservice

import (
	"dfs/internal/models"
	"sync"
	"time"
)

type MasterService struct {
	ChunkSize       uint64
	ChunkserverUrls map[string]time.Time
	Mutex           sync.Mutex
}

func NewMasterService(chunkSize uint64) *MasterService {
	return &MasterService{chunkSize, make(map[string]time.Time), sync.Mutex{}}
}

func (master *MasterService) HeartBeatListen(heartbeatMessage models.HeartbeatRequest) {
	master.Mutex.Lock()
	master.ChunkserverUrls[heartbeatMessage.Url] = time.Now()
	master.Mutex.Unlock()
}

func (master *MasterService) Upload(uploadReq models.UploadInitRequest) (any, error) {
	panic("not implemented")

}

func (master *MasterService) Get() (any, error) {
	panic("not implemented")

}

func (master *MasterService) View() (any, error) {
	panic("not implemented")
}

func (master *MasterService) SendUploadFileSuccessfully() (any, error) {
	panic("not implemented")
}
