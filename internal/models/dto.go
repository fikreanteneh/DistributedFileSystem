package models

type RegisterChunkServerRequest struct {
	Url string `json:"url"`
}

type RegisterChunkServerResponse struct {
}

type UploadInitRequest struct {
	FileName string `json:"fileName"`
	FileSize uint64 `json:"fileSize"`
}

type UploadInitResponse struct {
	Identifier     string `json:"identifier"`
	ChunkSize      uint64 `json:"chunkSize"`
	NumberOfChunks uint64 `json:"numberOfChunks"`
	ChunkServer    string `json:"chunkServer"`
}

type ChunkUploadSuccessRequest struct {
	ChunkIdentifier string `json:"chunkIdentifier"`
	ChunkServer     string `json:"chunkServer"`
}

type GetFileResponse struct {
	FileName  string
	Locations []string
}

type HeartbeatNotifier struct {
	Url  string `json:"url"`
	Port string `json:"port"`
	Ip   string `json:"ip"`
}

type ViewResponse struct {
	Files []FileMetadata `json:"files"`
}
