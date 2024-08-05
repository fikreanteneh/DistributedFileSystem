package chunkservice

type ChunkService struct {
}

func NewChunkService() *ChunkService {
	return &ChunkService{}
}

func (chunkService *ChunkService) Upload() (any, error) {
	panic("not implemented")
}

func (chunkService *ChunkService) Get() (any, error) {
	panic("not implemented")
}

func (chunkService *ChunkService) HeartBeatEmit() (any, error) {
	panic("not implemented")
}

func (chunkService *ChunkService) SendSuccessfulUpload() (any, error) {
	panic("not implemented")
}
