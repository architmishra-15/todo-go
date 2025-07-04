#include "parser.hpp"
#include <cctype>

namespace toml {

// ---------- Lexer Implementation ----------

Lexer::Lexer(const std::string& src)
    : source(src), pos(0), line(1) {}

char Lexer::advance() {
    if (pos < source.size()) {
        char c = source[pos++];
        if (c == '\n') ++line;
        return c;
    }
    return '\0';
}

char Lexer::peek() const {
    if (pos < source.size()) return source[pos];
    return '\0';
}

void Lexer::skipWhitespace() {
    while (true) {
        char c = peek();
        if (isspace(c)) {
            advance();
        } else if (c == '#') {
            // `#` will be considered as comments
            while (peek() != '\n' && peek() != '\0') advance();
        } else {
            break;
        }
    }
}

Token Lexer::lexString() {
    std::string val;
    // skip opening quote
    advance();
    while (peek() != '"' && peek() != '\0') {
        char c = advance();
        // basic escape sequences
        if (c == '\\' && peek() != '\0') {
            char escaped = advance();
            switch (escaped) {
                case 'n': val += '\n'; break;
                case 't': val += '\t'; break;
                case 'r': val += '\r'; break;
                case '\\': val += '\\'; break;
                case '"': val += '"'; break;
                default: val += escaped; break;
            }
        } else {
            val += c;
        }
    }
    if (peek() == '"') advance();
    return { TokenType::String, val, line };
}

Token Lexer::lexNumberOrBoolean() {
    std::string lit;
    bool isFloat = false;
    
    // negative numbers
    if (peek() == '-') {
        lit += advance();
    }
    
    while (isdigit(peek()) || peek() == '.') {
        if (peek() == '.') {
            if (isFloat) break; // Only one decimal point allowed
            isFloat = true;
        }
        lit += advance();
    }
    
    // Return appropriate numeric type
    if (isFloat) {
        return { TokenType::Float, lit, line };
    }
    return { TokenType::Integer, lit, line };
}

Token Lexer::lexIdentifier() {
    std::string id;
    while (isalnum(peek()) || peek() == '_' || peek() == '-') {
        id += advance();
    }
    
    // Check for boolean values
    if (id == "true" || id == "false") {
        return { TokenType::Boolean, id, line };
    }
    
    return { TokenType::Identifier, id, line };
}

Token Lexer::nextToken() {
    skipWhitespace();
    char c = peek();
    if (c == '\0') return { TokenType::EndOfFile, "", line };
    
    switch (c) {
        case '=': advance(); return { TokenType::Equal, "=", line };
        case '[': advance(); return { TokenType::LeftBracket, "[", line };
        case ']': advance(); return { TokenType::RightBracket, "]", line };
        case ',': advance(); return { TokenType::Comma, ",", line };
        case '.': advance(); return { TokenType::Dot, ".", line };
        case '"': return lexString();
        case '-':
            // Could be negative number or identifier
            if (isdigit(source[pos + 1])) {
                return lexNumberOrBoolean();
            }
            return lexIdentifier();
        default:
            if (isdigit(c)) return lexNumberOrBoolean();
            if (isalpha(c) || c == '_') return lexIdentifier();
            // unknown char
            advance();
            return { TokenType::Error, std::string(1, c), line };
    }
}

// ---------- Parser Implementation ----------

Parser::Parser(const std::vector<Token>& toks)
    : tokens(toks), current(0)
{
    rootTable = std::make_shared<Table>();
    currentTable = rootTable;
}

std::shared_ptr<Table> Parser::parse() {
    while (peek().type != TokenType::EndOfFile) {
        if (peek().type == TokenType::LeftBracket) {
            // table header: [section] or [section.subsection]
            advance(); // consume '['
            std::string tablePath = parseKey();
            consume(TokenType::RightBracket, "Expected ] after table name");
            
            // Handle nested tables ([something.example])
            currentTable = createNestedTable(rootTable, tablePath);
        } else {
            // key-value pair
            parseKeyValue(*currentTable);
        }
    }
    return rootTable;
}

Token Parser::peek() const {
    if (current < tokens.size()) {
        return tokens[current];
    }
    return { TokenType::EndOfFile, "", 0 };
}

Token Parser::advance() {
    if (current < tokens.size()) return tokens[current++];
    return { TokenType::EndOfFile, "", 0 };
}

bool Parser::match(TokenType expected) {
    if (peek().type == expected) {
        advance();
        return true;
    }
    return false;
}

void Parser::consume(TokenType expected, const std::string& errorMessage) {
    if (!match(expected)) {
        throw std::runtime_error(errorMessage + " at line " + std::to_string(peek().line));
    }
}

void Parser::parseKeyValue(Table& table) {
    std::string key = parseKey();
    consume(TokenType::Equal, "Expected = after key");
    Value val = parseValue();
    table.entries[key] = val;
}

std::string Parser::parseKey() {
    std::string key;
    Token t = advance();
    if (t.type != TokenType::Identifier) {
        throw std::runtime_error("Expected identifier for key at line " + std::to_string(t.line));
    }
    key = t.lexeme;
    
    // Handle dotted keys like [something.example]
    while (peek().type == TokenType::Dot) {
        advance(); // consume '.'
        Token next = advance();
        if (next.type != TokenType::Identifier) {
            throw std::runtime_error("Expected identifier after '.' in key at line " + std::to_string(next.line));
        }
        key += "." + next.lexeme;
    }
    
    return key;
}

Value Parser::parseValue() {
    Token t = peek();
    switch (t.type) {
        case TokenType::String:
            advance();
            return t.lexeme;
        case TokenType::Integer:
            advance();
            try {
                return std::stoll(t.lexeme);
            } catch (const std::exception&) {
                throw std::runtime_error("Invalid integer value: " + t.lexeme + " at line " + std::to_string(t.line));
            }
        case TokenType::Float:
            advance();
            try {
                return std::stold(t.lexeme);
            } catch (const std::exception&) {
                throw std::runtime_error("Invalid float value: " + t.lexeme + " at line " + std::to_string(t.line));
            }
        case TokenType::Boolean:
            advance();
            return (t.lexeme == "true");
        case TokenType::LeftBracket:
            return parseArray();
        default:
            throw std::runtime_error("Unexpected token '" + t.lexeme + "' in value at line " + std::to_string(t.line));
    }
}

std::vector<Value> Parser::parseArray() {
    std::vector<Value> elements;
    consume(TokenType::LeftBracket, "Expected [ to start array");
    
    // Handle empty array
    if (peek().type == TokenType::RightBracket) {
        advance();
        return elements;
    }
    
    while (peek().type != TokenType::RightBracket && peek().type != TokenType::EndOfFile) {
        elements.push_back(parseValue());
        if (!match(TokenType::Comma)) break;
    }
    consume(TokenType::RightBracket, "Expected ] to end array");
    return elements;
}

std::shared_ptr<Table> Parser::createNestedTable(std::shared_ptr<Table> root, const std::string& path) {
    std::vector<std::string> parts = splitPath(path);
    std::shared_ptr<Table> current = root;
    
    for (const auto& part : parts) {
        auto it = current->entries.find(part);
        if (it != current->entries.end()) {
            // Table already exists, use it
            if (std::holds_alternative<std::shared_ptr<Table>>(it->second)) {
                current = std::get<std::shared_ptr<Table>>(it->second);
            } else {
                throw std::runtime_error("Key '" + part + "' already exists as non-table value");
            }
        } else {
            // Create new table
            auto newTable = std::make_shared<Table>();
            current->entries[part] = newTable;
            current = newTable;
        }
    }
    
    return current;
}

std::vector<std::string> Parser::splitPath(const std::string& path) {
    std::vector<std::string> parts;
    std::string current;
    
    for (char c : path) {
        if (c == '.') {
            if (!current.empty()) {
                parts.push_back(current);
                current.clear();
            }
        } else {
            current += c;
        }
    }
    
    if (!current.empty()) {
        parts.push_back(current);
    }
    
    return parts;
}

// ---------- Convenience Function ----------

std::shared_ptr<Table> parseToml(const std::string& input) {
    Lexer lexer(input);
    std::vector<Token> tokens;
    Token tok;
    do {
        tok = lexer.nextToken();
        tokens.push_back(tok);
    } while (tok.type != TokenType::EndOfFile);

    Parser parser(tokens);
    return parser.parse();
}

} // namespace toml
