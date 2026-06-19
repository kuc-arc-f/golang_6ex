#include "processor.h"
#include <cstring>
#include <cstdlib>
#include <string>
#include <algorithm>
#include <sstream>
#include <cctype>
#include "my_rag.hpp"

extern "C" {

// malloc で確保したメモリを解放する
void free_string(char* str) {
    if (str != nullptr) {
        free(str);
    }
}    
/*
char* rag_add() {
    MyRag todo_helper("");
    todo_helper.rag_add_handler();
    std::string result = "OK";
    char* output = new char[result.length() + 1];
    strcpy(output, result.c_str());
    return output;    
}  
*/

char* rag_search(const char* input) {
    std::string input_str(input);
    MyRag todo_helper("");
    std::string result = todo_helper.rag_search_handler(input_str);
    char* output = new char[result.length() + 1];
    strcpy(output, result.c_str());
    return output;    
}

// Go から受け取った文字列を処理し、反転した結果と合わせて返す
char* process_string(const char* input) {
    std::string str(input);

    // 文字列を反転（UTF-8 バイト単位）
    std::string reversed(str.rbegin(), str.rend());

    // 結果を組み立て
    std::string result = "[C++ 処理結果] 受信: \"" + str + "\" → 反転: \"" + reversed + "\"";

    // malloc で確保して返す（Go 側で free_string を呼ぶ）
    char* output = (char*)malloc(result.length() + 1);
    std::strcpy(output, result.c_str());
    return output;
}

// 文字列を大文字に変換して返す（ASCII 部分のみ）
char* to_upper(const char* input) {
    std::string str(input);
    std::string upper_str;
    upper_str.reserve(str.size());

    for (char c : str) {
        upper_str += static_cast<char>(std::toupper(static_cast<unsigned char>(c)));
    }

    std::string result = "[C++ 大文字変換] \"" + std::string(input) + "\" → \"" + upper_str + "\"";

    char* output = (char*)malloc(result.length() + 1);
    std::strcpy(output, result.c_str());
    return output;
}

// 文字列の情報（バイト数、ASCII 文字数）を返す
char* string_info(const char* input) {
    std::string str(input);

    size_t byte_count = str.size();
    size_t ascii_count = 0;
    size_t non_ascii_count = 0;

    for (unsigned char c : str) {
        if (c < 128) {
            ascii_count++;
        } else {
            non_ascii_count++;
        }
    }

    std::ostringstream oss;
    oss << "[C++ 文字列情報] \"" << str << "\""
        << " | バイト数: " << byte_count
        << " | ASCII文字数: " << ascii_count
        << " | 非ASCII バイト数: " << non_ascii_count;

    std::string result = oss.str();
    char* output = (char*)malloc(result.length() + 1);
    std::strcpy(output, result.c_str());
    return output;
}


} // extern "C"
