package lexer

import (
	"testing"

	"github.com/naivary/codemark/internal/lexer/token"
)

// TODO: remove comma its unneccary
func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenOrder []Token
	}{
		{
			name:  "string",
			input: `+codemark:lexer:string="string"`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "string"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "bool without assignment",
			input: `+codemark:lexer:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "bool with assignment false",
			input: `+codemark:lexer:bool=false`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "false"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "bool with assignment true",
			input: `+codemark:lexer:bool=true`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "decimal",
			input: "+codemark:lexer:int=99",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "decimal negative sign",
			input: `+codemark:lexer:int=-99`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "-99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "decimal positve sign",
			input: `+codemark:lexer:int=+99`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "+99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "float",
			input: "+codemark:lexer:float=99.99",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:float"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.FLOAT, "99.99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "float negative sign",
			input: "+codemark:lexer:float=-99.99",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:float"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.FLOAT, "-99.99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "complex",
			input: "+codemark:lexer:complex=9i+9",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "9i+9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "complex wiht 0 imaginary",
			input: "+codemark:lexer:complex=0i+9",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "0i+9"),
				NewToken(token.EOF, ""),
			},
		},

		{
			name:  "complex negative sign",
			input: "+codemark:lexer:complex=-9i-9",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "-9i-9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "complex positive sign",
			input: "+codemark:lexer:complex=+9i+9",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:complex"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.COMPLEX, "+9i+9"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "list with all element types",
			input: `+codemark:lexer:list=[99, 99.99, 9i+9, true, false, "string"]`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:list"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.LBRACK, "["),
				NewToken(token.INT, "99"),
				NewToken(token.FLOAT, "99.99"),
				NewToken(token.COMPLEX, "9i+9"),
				NewToken(token.BOOL, "true"),
				NewToken(token.BOOL, "false"),
				NewToken(token.STRING, "string"),
				NewToken(token.RBRACK, "]"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "random input before",
			input: `Lorem ipsum dolor sit amet +jsonschema:validation:maximum=3`,
			tokenOrder: []Token{
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "string before with newline",
			input: `Lorem ipsum dolor sit amet 
					+codemark:lexer:string="string"`,
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
			input: `+codemark:lexer:int=99
					+codemark:lexer:string="string"
			`,
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
					+codemark:lexer:int=99
					+codemark:lexer:string="string"
					Lorem ipsum more documentation
			`,
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
			name:  "no identifier",
			input: `+=3`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "spaces before marker",
			input: `       +codemark:lexer:int=99`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "99"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "marker followed by text",
			input: `+codemark:lexer:bool=true Lorem`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:bool"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "unfinished assignment",
			input: `+codemark:lexer:assign.unfinished=`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:assign.unfinished"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "invalid string value",
			input: `+codemark:lexer:string="string\"`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, `string\"`),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier only domain",
			input: `+codemark="string"`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},

		{
			name:  "identifier only domain and resource",
			input: `+codemark:lexer="string"`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to start with a number",
			input: `+3codemark:lexer:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is allowed to include numbers",
			input: `+codemark_v1.1:lexer_v1.1:bool._required`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark_v1.1:lexer_v1.1:bool._required"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.BOOL, "true"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with underscore (first segment)",
			input: `+codemark_v1.1_:lexer:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with underscore (middle segment)",
			input: `+codemark:lexer_:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with underscore (end segment)",
			input: `+codemark:lexer:bool_`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with dot (first segment)",
			input: `+codemark.:lexer:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with dot (middle segment)",
			input: `+codemark:lexer.:bool`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "identifier is not allowed to end with dot (end segment)",
			input: `+codemark:lexer:bool.`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name: "newline followed by EOF",
			input: `+codemark:lexer:int=0x23f
`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:int"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.INT, "0x23f"),
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "multiline string",
			input: `+codemark:lexer:multiline.string=` + "`" + `string 
which is multi line` + "`",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:multiline.string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, `string 
which is multi line`),
				NewToken(token.EOF, ""),
			},
		},
		{
			name: "single line in multiline form",
			input: `+codemark:lexer:string="this is
multi line not allowed
`,
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, "this is"),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "multiline string with illegal character tick",
			input: `+codemark:lexer:multiline.string=` + "`" + `string which is multi` + "`" + `line` + "`",
			tokenOrder: []Token{
				NewToken(token.PLUS, "+"),
				NewToken(token.IDENT, "codemark:lexer:multiline.string"),
				NewToken(token.ASSIGN, "="),
				NewToken(token.STRING, `string which is multi`),
				NewToken(token.ERROR, ""),
			},
		},
		{
			name:  "documentation which contains marker",
			input: `this is a documentation string which contains a marker in the text +codemark:lexer:multiline.string=`,
			tokenOrder: []Token{
				NewToken(token.EOF, ""),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := Lex(tc.input)
			var i int
			for tk := range l.tokens {
				t.Log(tk.Kind)
				wantKind := tc.tokenOrder[i].Kind
				if tk.Kind != wantKind {
					t.Fatalf("kind's do not match. got: %s; want: %s", tk.Kind, wantKind)
				}
				if tk.Kind == token.ERROR || tk.Kind == token.EOF {
					t.Skipf("validation of ERROR or EOF tokens will be skipped")
					continue
				}
				wantValue := tc.tokenOrder[i].Value
				if tk.Value != wantValue {
					t.Fatalf("value's don't match. got: %s; want: %s", tk.Value, wantValue)
				}
				i++
			}
		})
	}
}
