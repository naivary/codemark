package main

import "testing"

func TestLexer_Lex(t *testing.T) {
    tests := []struct{
        name string
        input string
        isValid bool
        want bool
    }{
        {
            name:"",
            input: "dasdasd asd asdsadasd as +jsonschema:validation:maximum=3",
            isValid: true,
            want: true,
        },
        {
            name: "",
            input: "+jsonschema:validation:maximum=3",
            isValid: true,
            want: true,
        },
        {
            name: "",
            input: "+jsonschema:validation:required",
            isValid: true,
            want: true,
        },
        {
            name: "",
            input: "+jsonschema:validation:items=[value, 4, true]",
            isValid: true,
            want: true,
        },

    }

    for _, tc := range tests {
        l := lex(tc.name, tc.input)
        l.run()
        for token := range l.tokens {
            t.Log(token)
        }
    }
}
