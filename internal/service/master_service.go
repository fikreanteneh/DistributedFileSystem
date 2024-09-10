package service

import (
	"dfs/internal/interfaces"
	"dfs/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

type MasterService struct {
	repository         interfaces.FileRepositoryInterface
	chunkSize          uint64
	chunkServerDetails map[string]time.Time
	mutex              sync.Mutex
}

func NewMasterService(chunkSize uint64, repository interfaces.FileRepositoryInterface) *MasterService {
	return &MasterService{repository, chunkSize, make(map[string]time.Time), sync.Mutex{}}
}

func (master *MasterService) HeartBeatListen(heartbeatMessage *models.HeartbeatNotifier) {
	//TODO: Figure Out If the mutex is important
	// master.Mutex.Lock()
	master.chunkServerDetails[heartbeatMessage.Url] = time.Now()
	// master.Mutex.Unlock()
}

func (master *MasterService) Upload(uploadReq *models.UploadInitRequest) (*models.UploadInitResponse, error) {
	chunkSize := master.chunkSize
	numberOfChunks := uploadReq.FileSize / chunkSize
	identifier := uuid.New().String()

	keys := make([]string, 0, len(master.chunkServerDetails))
	for k := range master.chunkServerDetails {
		keys = append(keys, k)
	}

	res := models.UploadInitResponse{
		Identifier:     identifier,
		ChunkSize:      chunkSize,
		NumberOfChunks: numberOfChunks,
		ChunkServer:    keys[rand.Intn(len(keys))],
	}
	return &res, nil
}

func (master *MasterService) Get(fileId string) (*models.FileMetadata, error) {
	fileMetadata, err := master.repository.GetByID(fileId)
	if err != nil {
		return nil, err
	}
	return fileMetadata, nil
}

func (master *MasterService) View() (*[]*models.FileMetadata, error) {
	fileMetadata, err := master.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return fileMetadata, nil
}

func (master *MasterService) SendUploadFileSuccessfully() (any, error) {
	panic("not implemented")
}
