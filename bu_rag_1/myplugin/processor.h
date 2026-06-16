#ifndef PROCESSOR_H
#define PROCESSOR_H

#ifdef __cplusplus
extern "C" {
#endif

// 返却された文字列のメモリを解放する
void free_string(char* str);

char* rag_search(const char* input);

// Go から受け取った文字列を処理し、結果を返す
char* process_string(const char* input);

// 文字列を大文字に変換して返す
char* to_upper(const char* input);

// 文字列の情報（長さ、バイト数など）を返す
char* string_info(const char* input);


#ifdef __cplusplus
}
#endif

#endif // PROCESSOR_H
