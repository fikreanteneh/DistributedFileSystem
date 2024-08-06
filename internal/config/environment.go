package config

type Environment struct {
	MasterServerURL  string
	ChunkServerURL   string
	ChunkServerDIR   string
	MasterServerPort string
	ChunkServerPort  string
}

func NewEnvironment() *Environment {
	return &Environment{
		MasterServerURL:  "http://localhost:8000/",
		MasterServerPort: "8000",
		ChunkServerURL:   "http://localhost:8001/",
		ChunkServerPort:  "8001",
		ChunkServerDIR:   "/data",
	}
}
