package lexer

import (
	"testing"

	"github.com/naivary/codemark/internal/lexer/token"
)

// TODO: better testing with expected token order
// TODO: remove comma its unneccary
func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		isValid    bool
		tokenOrder []Token
	}{
		{
			name:    "string",
			input:   `+codemark:lexer:string="string"`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "string"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "bool without assignment",
			input:   `+codemark:lexer:bool`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "bool with assignment false",
			input:   `+codemark:lexer:bool=false`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "false"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "bool with assignment true",
			input:   `+codemark:lexer:bool=true`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "decimal",
			input:   "+codemark:lexer:int=99",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "decimal negative sign",
			input:   `+codemark:lexer:int=-99`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "-99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "decimal positve sign",
			input:   `+codemark:lexer:int=+99`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "float",
			input:   "+codemark:lexer:float=99.99",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:float"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.FLOAT, "99.99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "float negative sign",
			input:   "+codemark:lexer:float=-99.99",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:float"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.FLOAT, "-99.99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "complex",
			input:   "+codemark:lexer:complex=9i+9",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "9i+9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "complex negative sign",
			input:   "+codemark:lexer:complex=-9i-9",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "-9i-9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "complex positive sign",
			input:   "+codemark:lexer:complex=+9i+9",
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "+9i+9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:    "list with all element types",
			input:   `+codemark:lexer:list=[99, 99.99, 9i+9, true, false, "string"]`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:list"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.LBRACK, "["),
				NewToken(token.INT, "99"),
				NewToken(token.COMMA, ","),
				NewToken(token.FLOAT, "99.99"),
				NewToken(token.COMMA, ","),
				NewToken(token.COMPLEX, "9i+9"),
				NewToken(token.COMMA, ","),
				NewToken(token.BOOL, "true"),
				NewToken(token.COMMA, ","),
				NewToken(token.BOOL, "false"),
				NewToken(token.COMMA, ","),
				NewToken(token.STRING, "string"),
				NewToken(token.RBRACK, "]"),
				NewToken(token.EOF, ""),
			},
		},
		// TODO: remove duplicated test and take over good casses
		{
			name:    "random input before",
			input:   `Lorem ipsum dolor sit amet +jsonschema:validation:maximum=3`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "string before with newline",
			input: `Lorem ipsum dolor sit amet 
					+codemark:lexer:string="string"`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "string"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "multi marker without doc",
			input: `+codemark:lexer:string="string"
					+codemark:lexer:int=99
			`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "string"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "multi marker with doc",
			input: `Lorem ipsum documentation
					+codemark:lexer:string="string"
					+codemark:lexer:int=99
					Lorem ipsum more documentation
			`,
			isValid: true,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "string"),
				NewToken(token.EOF, ""),
			},
		},

		{
			name:    "no identifier",
			input:   `+=3`,
			isValid: false,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:    "spaces before marker",
			input:   `       +codemark:lexer:int=99`,
			isValid: true,
		},
		{
			name:    "multi marker without doc",
			input:   `+codemark:lexer:bool=true Lorem`,
			isValid: false,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "true"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:    "unfinished assignment",
			input:   `+codemark:lexer:assign.unfinished=`,
			isValid: false,
		},
		{
			name:    "invalid string value",
			input:   `+codemar:lexer:string="string\"`,
			isValid: false,
		},
		{
			name:    "identifier only domain",
			input:   `+codemark="string"`,
			isValid: false,
		},

		{
			name:    "identifier only domain and resource",
			input:   `+codemark:lexer="string"`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to start with a number",
			input:   `+3codemark:lexer:bool`,
			isValid: false,
		},
		{
			name:    "identifier is allowed to include numbers",
			input:   `+codemark_v1.1:lexer_v1.1:bool._required`,
			isValid: true,
		},
		{
			name:    "identifier is not allowed to end with underscore (first segment)",
			input:   `+codemark_v1.1_:lexer:bool`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with underscore (middle segment)",
			input:   `+codemark:lexer_:bool`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with underscore (end segment)",
			input:   `+codemark:lexer:bool_`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (first segment)",
			input:   `+codemark.:lexer:bool`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (middle segment)",
			input:   `+codemark:lexer.:bool`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (end segment)",
			input:   `+codemark:lexer:bool.`,
			isValid: false,
		},
		{
			name: "newline followed by EOF",
			input: `+codemark:lexer:int=0x23f
`,
			isValid: true,
		},
		{
			name: "multiline string",
			input: `+codemark:lexer:multiline.string=` + "`" + `string 
which is multi line` + "`",
			isValid: true,
		},
		{
			name: "single line in multile form",
			input: `+codemark:lexer:string="this is
multi line not allowed
`,
			isValid: false,
		},
		{
			name: "multiline string with illegal character tick",
			input: `+codemark:lexer:multiline.string=` + "`" + `string 
which is multi` + "`" + `line` + "`",
			isValid: false,
		},
		{
			name:    "documentation which contains marker",
			input:   `this is a documentation string which contains a marker in the text +codemark:lexer:multiline.string=`,
			isValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := Lex(tc.input)
			for tk := range l.tokens {
				t.Log(tk)
				if tk.Kind == token.ERROR && tc.isValid {
					t.Fatalf("expected to lex correctly: `%s`. Error is: %s", tc.input, tk.Value)
				}
				if tk.Kind == token.EOF && !tc.isValid {
					t.Fatalf("expected an error and didn't get one for `%s`", tc.input)
				}
			}
		})
	}
}
