# bu_rag_2

 Version: 0.9.1

 Author  :

 date    : 2026/06/19

 update :

***

Golang Bubble Tea , RAG Qdrant

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

https://github.com/kuc-arc-f/cpp_3ex/tree/main/rag_19

***
* RAG APP

![img1](/images/bu_rag_2.png)

***
### related

https://github.com/charmbracelet/bubbletea

***
* LIB add

```
sudo apt install uuid-dev
sudo apt install libcurl4-openssl-dev
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
* env value
* GEMINI_API_KEY SET
```
export GEMINI_API_KEY=your-key
```

***
* C++ build
```
make plugin
```

* build
```
go mod init example.com/bu-rag-2
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

