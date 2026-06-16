// Package main は Go plugin として C++ の文字列処理関数をエクスポートする。
// go build -buildmode=plugin でビルドし、メインプログラムから plugin パッケージで読み込む。
package main

/*
#cgo CXXFLAGS: -std=c++17 -O2
#cgo LDFLAGS: -lstdc++ -lsqlite3 -lcurl -luuid
#include "processor.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func TodoAdd(msg string) string {
    cMsg := C.CString(msg)
    defer C.free(unsafe.Pointer(cMsg))

    cResult := C.todo_add(cMsg)
    defer C.free_string(cResult)

    return C.GoString(cResult)
}
func TodoList() string {
    //cMsg := C.CString(msg)
    //defer C.free(unsafe.Pointer(cMsg))
    cResult := C.todo_list()
    defer C.free_string(cResult)

    return C.GoString(cResult)
}

/*
func RagAdd() string {
    cResult := C.rag_add()
    defer C.free_string(cResult)

    return C.GoString(cResult)
}


func RagSearch(msg string) string {
    cMsg := C.CString(msg)
    defer C.free(unsafe.Pointer(cMsg))

    cResult := C.rag_search(cMsg)
    defer C.free_string(cResult)

    return C.GoString(cResult)
}
*/

// ProcessString は Go の文字列を C++ に送信し、処理結果を受け取って返す。
// C++ 側で文字列の反転処理が行われる。
func ProcessString(msg string) string {
    cMsg := C.CString(msg)
    defer C.free(unsafe.Pointer(cMsg))

    cResult := C.process_string(cMsg)
    defer C.free_string(cResult)

    return C.GoString(cResult)
}

// ToUpper は Go の文字列を C++ に送信し、大文字変換結果を受け取って返す。
func ToUpper(msg string) string {
    cMsg := C.CString(msg)
    defer C.free(unsafe.Pointer(cMsg))

    cResult := C.to_upper(cMsg)
    defer C.free_string(cResult)

    return C.GoString(cResult)
}

// StringInfo は Go の文字列を C++ に送信し、文字列情報を受け取って返す。
func StringInfo(msg string) string {
    cMsg := C.CString(msg)
    defer C.free(unsafe.Pointer(cMsg))

    cResult := C.string_info(cMsg)
    defer C.free_string(cResult)

    return C.GoString(cResult)
}
