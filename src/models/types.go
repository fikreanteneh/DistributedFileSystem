package models



type RegisterChunkserverRequest struct {
	Url string `json:"url"`
}

type RegisterChunkserverResponse struct {
}

type UploadInitRequest struct {
	FileName string `json:"fileName"`
	FileSize uint64 `json:"fileSize"`
}

type UploadInitResponse struct {
	Identifier     string   `json:"identifier"`
	ChunkSize      uint64   `json:"chunkSize"`
	NumberOfChunks uint64   `json:"numberOfChunks"`
	Chunkservers   []string `json:"chunkservers"`
}

type ChunkUploadSuccessRequest struct {
	ChunkIdentifier string `json:"chunkIdentifier"`
	Chunkserver     string `json:"chunkserver"`
}

type GetResponse struct {
	FileName  string
	Locations []string
}


type HeartbeatRequest struct {
	Url string `json:"url"`
	Port string `json:"port"`
	Ip string `json:"ip"`
}

type ViewResponse struct {
	Files []FileMetadata `json:"files"`
}