package main

import "testing"

func TestRegistry_Lookup(t *testing.T) {
    reg := NewRegistry()

    type args struct {
        name string
        want *Definition
    }

    defs := map[string]*Definition {
        "test:path:mark": &Definition{
            Name: "test:path:mark",
            TargetType: TargetField,
        },
    }

    tests := []args{
        {
            name: "test:path:mark",
            want: defs["test:path:mark"],
        },
    }

    for _, def := range defs {
        if err := reg.Register(def); err != nil{
            t.Error(err)
        }
    }

    for _, tc := range tests {
        def := reg.Lookup(tc.name)
        if def != tc.want {
            t.Errorf("definition is not as expected. Got: %v. Wanted: %v", def, tc.want)
        }
    }
}
