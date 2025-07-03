#pragma once

#include <memory>
#include <string>
#include <unordered_map>
#include <variant>
#include <vector>
#include <stdexcept>

namespace toml {

enum class TokenType {
    Identifier,   
    Equal,
    String,
    Integer,
    Float,
    Boolean,
    Comma,
    LeftBracket,
    RightBracket,
    EndOfFile,
    Error
};

// Token representation
struct Token {
    TokenType type;
    std::string lexeme;
    int line;
};

// Variant type to hold parsed TOML values
using Integer = long long;
using Float   = long double;
using Boolean = bool;
using String  = std::string;

// Forward declare Table
struct Table;

using Value = std::variant<
    Integer,
    Float,
    Boolean,
    String,
    std::vector<Value>,
    std::shared_ptr<Table>
>;

// mapping of string keys to Values
struct Table {
    std::unordered_map<std::string, Value> entries;
};

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

class Parser {
public:
    explicit Parser(const std::vector<Token>& tokens);
    std::shared_ptr<Table> parse();

private:
    const std::vector<Token> tokens;
    size_t current = 0;

    Token peek() const;
    Token advance();
    bool match(TokenType expected);
    void consume(TokenType expected, const std::string& errorMessage);

    void parseKeyValue(Table& table);
    std::string parseKey();
    Value parseValue();
    std::vector<Value> parseArray();

    std::shared_ptr<Table> currentTable;
};

std::shared_ptr<Table> parseToml(const std::string& input);

} // namespace toml

