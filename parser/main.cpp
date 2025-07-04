#include "parser.hpp"
#include <fstream>
#include <iostream>
#include <memory>

// Function to print TOML structure hierarchically
void printValue(const toml::Value& value, int indent = 0) {
    std::string indentStr(indent, ' ');
    
    if (std::holds_alternative<toml::Integer>(value)) {
        std::cout << std::get<toml::Integer>(value);
    } else if (std::holds_alternative<toml::Float>(value)) {
        std::cout << std::get<toml::Float>(value);
    } else if (std::holds_alternative<toml::Boolean>(value)) {
        std::cout << (std::get<toml::Boolean>(value) ? "true" : "false");
    } else if (std::holds_alternative<toml::String>(value)) {
        std::cout << "\"" << std::get<toml::String>(value) << "\"";
    } else if (std::holds_alternative<toml::Array>(value)) {
        const auto& arr = std::get<toml::Array>(value);
        std::cout << "[";
        for (size_t i = 0; i < arr.size(); ++i) {
            if (i > 0) std::cout << ", ";
            printValue(arr[i], indent);
        }
        std::cout << "]";
    } else if (std::holds_alternative<std::shared_ptr<toml::Table>>(value)) {
        const auto& table = std::get<std::shared_ptr<toml::Table>>(value);
        std::cout << "{\n";
        for (const auto& [key, val] : table->entries) {
            std::cout << indentStr << "  " << key << " = ";
            printValue(val, indent + 2);
            std::cout << "\n";
        }
        std::cout << indentStr << "}";
    }
}

void printTable(const std::shared_ptr<toml::Table>& table, const std::string& prefix = "", int indent = 0) {
    std::string indentStr(indent, ' ');
    
    for (const auto& [key, value] : table->entries) {
        std::string fullKey = prefix.empty() ? key : prefix + "." + key;
        
        if (std::holds_alternative<std::shared_ptr<toml::Table>>(value)) {
            const auto& subTable = std::get<std::shared_ptr<toml::Table>>(value);
            std::cout << indentStr << "[" << fullKey << "] -> Table\n";
            printTable(subTable, fullKey, indent + 2);
        } else {
            std::cout << indentStr << fullKey << " = ";
            printValue(value, indent);
            std::cout << " (";
            
            // Print type
            if (std::holds_alternative<toml::Integer>(value)) {
                std::cout << "Integer";
            } else if (std::holds_alternative<toml::Float>(value)) {
                std::cout << "Float";
            } else if (std::holds_alternative<toml::Boolean>(value)) {
                std::cout << "Boolean";
            } else if (std::holds_alternative<toml::String>(value)) {
                std::cout << "String";
            } else if (std::holds_alternative<toml::Array>(value)) {
                std::cout << "Array";
            }
            
            std::cout << ")\n";
        }
    }
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: parser.exe <config.toml>\n";
        return 1;
    }
    
    // Read entire TOML file into a string
    std::ifstream infile(argv[1]);
    if (!infile) {
        std::cerr << "Failed to open file: " << argv[1] << "\n";
        return 1;
    }
    std::string contents((std::istreambuf_iterator<char>(infile)),
                         std::istreambuf_iterator<char>());

    try {
        // Parse TOML into a root table
        std::shared_ptr<toml::Table> root = toml::parseToml(contents);

        // Print hierarchical structure
        std::cout << "Parsed TOML structure:\n";
        std::cout << "====================\n";
        printTable(root);
        
    } catch (const std::exception& ex) {
        std::cerr << "Error parsing TOML: " << ex.what() << "\n";
        return 1;
    }

    return 0;
}
