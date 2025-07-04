#pragma once

#include <memory>
#include <string>
#include <unordered_map>
#include <variant>
#include <vector>
#include <stdexcept>

namespace toml {

// Supported token types in our minimal parser
enum class TokenType {
    Identifier,   // key names
    Equal,        // '='
    String,       // "..."
    Integer,      // 123
    Float,        // 123.45
    Boolean,      // true/false
    Comma,        // ','
    LeftBracket,  // '['
    RightBracket, // ']'
    Dot,          // '.'
    EndOfFile,
    Error
};

// Token representation
struct Token {
    TokenType type;
    std::string lexeme;
    int line;
};

// Type aliases for TOML values
using Integer = long long;
using Float   = long double;
using Boolean = bool;
using String  = std::string;

// Forward declarations
struct Table;
struct Value;

// A TOML value can be one of: scalar, array, or table
using Array = std::vector<Value>;

struct Value : std::variant<
    Integer,
    Float,
    Boolean,
    String,
    Array,
    std::shared_ptr<Table>
> {
    using variant::variant;  // inherit constructors
};

// A TOML table: mapping of string keys to Values
struct Table {
    std::unordered_map<std::string, Value> entries;
};

// Lexer: transforms raw TOML text into tokens
class Lexer {
public:
    explicit Lexer(const std::string& source);
    Token nextToken();

private:
    const std::string source;
    size_t pos = 0;
    int line = 1;

    char advance();
    char peek() const;
    void skipWhitespace();
    Token lexString();
    Token lexNumberOrBoolean();
    Token lexIdentifier();
};

// Parser: builds a Table from a token stream
class Parser {
public:
    explicit Parser(const std::vector<Token>& tokens);
    std::shared_ptr<Table> parse();

private:
    const std::vector<Token> tokens;
    size_t current = 0;
    
    std::shared_ptr<Table> rootTable;
    std::shared_ptr<Table> currentTable;

    Token peek() const;
    Token advance();
    bool match(TokenType expected);
    void consume(TokenType expected, const std::string& errorMessage);

    void parseKeyValue(Table& table);
    std::string parseKey();
    Value parseValue();
    Array parseArray();
    
    // Helper methods for nested tables
    std::shared_ptr<Table> createNestedTable(std::shared_ptr<Table> root, const std::string& path);
    std::vector<std::string> splitPath(const std::string& path);
};

// Convenience function: parse TOML text into a root Table
std::shared_ptr<Table> parseToml(const std::string& input);

} // namespace toml
