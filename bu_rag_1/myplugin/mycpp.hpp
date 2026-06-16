#pragma once
#include <iostream>
#include <string>

class Printer {
public:
    // C++本来のstd::stringを使うメソッド
    void PrintMessage(const std::string& msg){
        std::cout << "[C++側が出力] " << msg << std::endl;

    }
};

class Provider {
public:
    std::string GetMessage() {
        return "from C++の世界からこんにちは！";
    }
};
