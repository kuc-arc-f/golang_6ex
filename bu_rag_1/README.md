# bu_rag_1

 Version: 0.9.1

 Author  :

 date    : 2026/06/16

 update :

***

Golang Bubble Tea , RAG SQLite

* embedding : Gemini-embedding-001
* model: gemma-4-E2B
* llama.cpp , llama-server
* go version go1.26.1 linux/amd64
* gcc version 14.2.0 
* make
* Golnag C++ call , Shared Object
* Linux

***
### vector data add

https://github.com/kuc-arc-f/golang_cpp1/tree/main/go_rag_1

***
* RAG APP

![img1](/images/bu_rag_1.png)


***
### related

https://github.com/charmbracelet/bubbletea

***
* LIB add
```
sudo apt install libcurl4-openssl-dev
sudo apt-get install libsqlite3-dev
sudo apt install nlohmann-json3-dev
```
***
* llama-server start
* port 8090: gemma-4-E2B

```
#gemma-4-E2B

/usr/local/llama-b8642/llama-server -m /var/lm_data/unsloth/gemma-4-E2B-it-Q4_K_S.gguf \
 --chat-template-kwargs '{"enable_thinking": false}' --port 8090 
```

***
* db file: example.db copy
* golang_cpp1/go_rag_1 , from folder

***
* env value
* GEMINI_API_KEY SET
```
export GEMINI_API_KEY=
```

***
* C++ build
```
make plugin
```

* build
```
go mod init example.com/bu-rag-1
go mod tidy

go build main.go
```
***
* start
```
./main
```
***
* operate
* text input , Enter key

***
### blog

https://zenn.dev/knaka0209/scraps/ff43e6d8092886


