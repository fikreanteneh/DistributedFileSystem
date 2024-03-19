#!/bin/bash

go run src/*.go master --port 8000 --chunkSize 1000 &
go run src/*.go chunkserver --master http://localhost:8000/ --port 8001 --dir ./data1 &
go run src/*.go chunkserver --master http://localhost:8000/ --port 8002 --dir ./data2 &
