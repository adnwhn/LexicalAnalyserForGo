package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Tipul TokenType este un alias pentru un sir de caractere care reprezinta tipul unui token
type TokenType string

// Definire tipuri de tokeni posibili
const (
	TokenIdent     TokenType = "Identificator"
	TokenKeyword   TokenType = "Cuvant cheie"
	TokenInt       TokenType = "Numar intreg"
	TokenFloat     TokenType = "Numar zecimal"
	TokenDelimiter TokenType = "Delimitator"
	TokenString    TokenType = "String"
	TokenComment   TokenType = "Comentariu"
	TokenOperator  TokenType = "Operator"
	TokenEOF       TokenType = "Sfarsit de fisier"
	TokenError     TokenType = "Eroare lexicala"
)

// Structura pentru un token cu informatii despre tipul sau, lungime, linie, pozitie in sir, valoare
type Token struct {
	Type   TokenType
	Length int
	Line   int
	Pos    int
	Value  string
}

// Map care contine cuvintele cheie si valorile lor asociate ca tipuri de tokeni
var keywords = map[string]TokenType{
	"break":       TokenKeyword,
	"case":        TokenKeyword,
	"chan":        TokenKeyword,
	"const":       TokenKeyword,
	"continue":    TokenKeyword,
	"default":     TokenKeyword,
	"defer":       TokenKeyword,
	"else":        TokenKeyword,
	"fallthrough": TokenKeyword,
	"for":         TokenKeyword,
	"func":        TokenKeyword,
	"go":          TokenKeyword,
	"goto":        TokenKeyword,
	"if":          TokenKeyword,
	"import":      TokenKeyword,
	"interface":   TokenKeyword,
	"map":         TokenKeyword,
	"package":     TokenKeyword,
	"range":       TokenKeyword,
	"return":      TokenKeyword,
	"select":      TokenKeyword,
	"struct":      TokenKeyword,
	"switch":      TokenKeyword,
	"type":        TokenKeyword,
	"var":         TokenKeyword,
}

// Slice care contine toti operatorii in limbajul Go
var operators = []string{
	":", "+", "-", "=", "*", "/", "%", "<", ">", "!", "&", "|", "^",
	"++", "--", "==", "<<", ">>", "&&", "||",
	":=", "+=", "-=", "*=", "/=", "%=", "<=", ">=", "!=", "&=", "|=", "^=", "<-", "->",
}

// Functia care verifica daca elementul curent apartine unui slice
func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// Functia care analizeaza si returneaza informatiile despre tokenul curent sau semnaleaza o eroare lexicala
func scanToken(input string, pos int, line int) (token Token, newPos int, newLine int, err error) {
	// Sare peste spatii, tab-uri si linii noi
	for pos < len(input) && (unicode.IsSpace(rune(input[pos])) || input[pos] == '\t' || (input[pos] == '\r' && input[pos+1] == '\n')) {
		if input[pos] == '\r' && input[pos+1] == '\n' {
			line++
		}
		pos++
	}

	// Daca a ajuns la sfarsitul fisierului, returneaza un token de tip EOF
	if pos >= len(input) {
		return Token{Type: TokenEOF, Length: 0, Line: line, Pos: pos}, pos, line, nil
	}

	// Determina tipul tokenului si lungimea sa
	switch {
	// Identificator
	case unicode.IsLetter(rune(input[pos])):
		start := pos
		for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos]))) {
			pos++
		}
		identifier := input[start:pos]
		// Keyword
		if keywordType, isKeyword := keywords[identifier]; isKeyword {
			return Token{Type: keywordType, Length: pos - start, Line: line, Pos: start, Value: identifier}, pos, line, nil
		}
		return Token{Type: TokenIdent, Length: pos - start, Line: line, Pos: start, Value: identifier}, pos, line, nil
	case input[pos] == '_' && unicode.IsLetter(rune(input[pos + 1])):
		start := pos
		for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos]))) || (input[pos] == '_') {
			pos++
		}
		return Token{Type: TokenIdent, Length: pos - start, Line: line, Pos: start, Value: input[start:pos]}, pos, line, nil
	// Int
	case unicode.IsDigit(rune(input[pos])):
		start := pos
		for pos < len(input) && (unicode.IsDigit(rune(input[pos])) || input[pos] == '.') {
			pos++
		}
		numberFound := input[start:pos]
		// Float
		if strings.Contains(input[start:pos], ".") {
			return Token{Type: TokenFloat, Length: pos - start, Line: line, Pos: start, Value: numberFound}, pos, line, nil
		}
		return Token{Type: TokenInt, Length: pos - start, Line: line, Pos: start, Value: numberFound}, pos, line, nil
	// Float subunitar
	case input[pos] == '.' && unicode.IsDigit(rune(input[pos+1])):
		start := pos
		for pos < len(input) && (unicode.IsDigit(rune(input[pos+1]))) {
			pos++
		}
		subunit := input[start : pos+1]
		return Token{Type: TokenFloat, Length: pos - start + 1, Line: line, Pos: start, Value: subunit}, pos + 1, line, nil
	// Delimitator
	case input[pos] == '(' || input[pos] == ')' || input[pos] == '[' || input[pos] == ']' || input[pos] == '{' || input[pos] == '}' || 
		input[pos] == ',' || input[pos] == ';' || (input[pos] == '.' && unicode.IsLetter(rune(input[pos-1])) && unicode.IsLetter(rune(input[pos+1]))):
		// Verificare inchidere delimitator (), [], {}
		if input[pos] == '(' {
			start := pos
			pos++
			for pos < len(input) && input[pos] != ')' {
				if input[pos] == ')' {
					break
				}
				pos++
			}
			if pos >= len(input) {
				return Token{Type: TokenError, Length: pos - start, Line: line, Pos: start, Value: input[start:pos]}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, pozitia %d: delimitatorul %s nu este inchis", line, start, string(input[start]))
			}
			return Token{Type: TokenDelimiter, Length: 1, Line: line, Pos: start, Value: string(input[start])}, start + 1, line, nil
		} else if input[pos] == '[' {
			start := pos
			pos++
			for pos < len(input) && input[pos] != ']' {
				if input[pos] == ']' {
					break
				}
				pos++
			}
			if pos >= len(input) {
				return Token{Type: TokenError, Length: pos - start, Line: line, Pos: start, Value: input[start:pos]}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, pozitia %d: delimitatorul %s nu este inchis", line, start, string(input[start]))
			}
			return Token{Type: TokenDelimiter, Length: 1, Line: line, Pos: start, Value: string(input[start])}, start + 1, line, nil
		} else if input[pos] == '{' {
			start := pos
			pos++
			for pos < len(input) && input[pos] != '}' {
				if input[pos] == '}' {
					break
				}
				pos++
			}
			if pos >= len(input) {
				return Token{Type: TokenError, Length: pos - start, Line: line, Pos: start, Value: input[start:pos]}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, pozitia %d: delimitatorul %s nu este inchis", line, start, string(input[start]))
			}
			return Token{Type: TokenDelimiter, Length: 1, Line: line, Pos: start, Value: string(input[start])}, start + 1, line, nil
		}
		return Token{Type: TokenDelimiter, Length: 1, Line: line, Pos: pos, Value: string(input[pos])}, pos + 1, line, nil
	// String
	case input[pos] == '"':
		start := pos
		pos++
		for pos < len(input) && input[pos] != '"' {
			if input[pos] == '"' {
				break
			}
			pos++
		}
		if pos >= len(input) {
			return Token{Type: TokenError, Length: pos - start, Line: line, Pos: start, Value: input[start:pos]}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, pozitia %d: string neterminat", line, start)
		}
		stringFound := input[start+1 : pos]
		return Token{Type: TokenString, Length: pos - start + 1 - 2, Line: line, Pos: start, Value: stringFound}, pos + 1, line, nil
	// Comentariu + Operator
	case Contains(operators, string(input[pos])):
		// Comentarii
		if input[pos] == '/' && (input[pos+1] == '/' || input[pos+1] == '*') {
			// Single-line
			if input[pos+1] == '/' {
				start := pos + 2
				pos = start + 1
				for pos < len(input) && !(input[pos] == '\r' && input[pos+1] == '\n') {
					pos++
				}
				comment := input[start:pos]
				return Token{Type: TokenComment, Length: pos - start, Line: line, Pos: start, Value: comment}, pos, line, nil
			}
			// Multi-line
			if input[pos+1] == '*' {
				start := pos + 2
				pos = start + 1
				for pos < len(input) && !(input[pos] == '*' && input[pos+1] == '/') {
					if input[pos] == '*' && input[pos+1] == '/' {
						break
					}
					pos++
				}
				if pos >= len(input) {
					return Token{Type: TokenError, Length: pos - start, Line: line, Pos: start, Value: input[start-2 : pos]}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, pozitia %d: comentariul nu este inchis", line, start)
				}
				comment := input[start:pos]
				return Token{Type: TokenComment, Length: pos - start + 1 - 2, Line: line, Pos: start, Value: comment}, pos + 2, line, nil
			}
		}
		// Operator
		if Contains(operators, string(input[pos+1])) {
			return Token{Type: TokenOperator, Length: 2, Line: line, Pos: pos, Value: string(input[pos : pos+2])}, pos + 2, line, nil
		} else {
			return Token{Type: TokenOperator, Length: 1, Line: line, Pos: pos, Value: string(input[pos])}, pos + 1, line, nil
		}
	default:
		// Daca e intalnit un caracter care nu apartine limbajului, e semnalata eroarea lexicala si se continua analiza
		return Token{Type: TokenError, Length: 1, Line: line, Pos: pos, Value: string(input[pos])}, pos + 1, line, fmt.Errorf("! Eroare lexicala la linia %d, poziția %d: caracterul '%c' nu apartine limbajului Go", line, pos, input[pos])
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Missing text file!")

	} else {
		input, _ := os.ReadFile(arguments[1])

		// Initializare variabile ce tin cont de pozitia si linia curenta
		pos := 0
		line := 1

		// Parcurgere input, se scaneaza tokenii pana se ajunge la sfarsitul fisierului
		for pos < len(input) {
			// Scanare token curent
			token, newPos, newLine, err := scanToken(string(input), pos, line)
			if err != nil {
				fmt.Println(err)
			}

			// Actualizare pozitie
			pos = newPos
			line = newLine

			// Terminare dupa ce se ajunge la sfarsitul fisierului
			if token.Type == TokenEOF {
				break
			}

			// Afisare token curent
			fmt.Printf("Token: '%s', Tip: %s, Lungime: %d, Linie: %d, Pozitie: %d\n", token.Value, token.Type, token.Length, token.Line, token.Pos)
		}
	}

}
