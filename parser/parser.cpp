#include "parser.hpp"
#include <cctypes>

namespace toml {

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
      }
      else if ( c == '#' ) {
        while (peek() != '\n' && peek() != '\0') advance();
      }
      else {
        break;
      }
    }
  }
} // namespace toml
