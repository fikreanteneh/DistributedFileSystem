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

func (master *MasterService) HeartBeatListen(heartbeatMessage *models.HeartbeatNotifier) {
	master.Mutex.Lock()
	master.ChunkserverUrls[heartbeatMessage.Url] = time.Now()
	master.Mutex.Unlock()
}

func (master *MasterService) Upload(uploadReq *models.UploadInitRequest) (*models.UploadInitRequest, error) {
	panic("not implemented")

}

func (master *MasterService) Get(fileId string) (*models.FileMetadata, error) {
	panic("not implemented")

}

func (master *MasterService) View() (*[]models.FileMetadata, error) {
	panic("not implemented")
}

func (master *MasterService) SendUploadFileSuccessfully() (any, error) {
	panic("not implemented")
}
