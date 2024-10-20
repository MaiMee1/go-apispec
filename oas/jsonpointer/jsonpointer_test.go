package jsonpointer

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPtr_Tokens(t *testing.T) {
	var testCases = []struct {
		pointer Ptr
		len     int
	}{
		{"", 0},
		{"/foo", 1},
		{"/foo/0", 2},
		{"/", 1},
		{"/a/b", 2},
		{"/c%d", 1},
		{"/e^f", 1},
		{"/g|h", 1},
		{"/\\j", 1},
		{"/\"l", 1},
		{"/ ", 1},
		{"/m~0n", 1},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tokens, err := tt.pointer.tokens()
			if err != nil {
				t.Errorf("got error %v, want <nil>", err)
			}
			if len(tokens) != tt.len {
				t.Errorf("got %d tokens, want %d", len(tokens), tt.len)
			}
		})
	}
}

func TestPtr_Access(t *testing.T) {
	var doc interface{}
	if err := json.Unmarshal([]byte(`{"foo":["bar","baz"],"":0,"a/b":1,"c%d":2,"e^f":3,"g|h":4,"i\\j":5,"k\"l":6," ":7,"m~n":8}`), &doc); err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		document interface{}
		pointer  Ptr
		want     interface{}
		err      error
	}{
		{doc, "", doc, nil},
		{doc, "/foo", []string{"bar", "baz"}, nil},
		{doc, "/foo/0", "bar", nil},
		{doc, "/", 0, nil},
		{doc, "/a~1b", 1, nil},
		{doc, "/c%d", 2, nil},
		{doc, "/e^f", 3, nil},
		{doc, "/g|h", 4, nil},
		{doc, "/i\\j", 5, nil},
		{doc, "/k\"l", 6, nil},
		{doc, "/ ", 7, nil},
		{doc, "/m~0n", 8, nil},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v, err := tt.pointer.Access(doc)
			if fmt.Sprint(err) != fmt.Sprint(tt.err) {
				t.Errorf("got error %v, want %v", err, tt.err)
			}
			if fmt.Sprint(v) != fmt.Sprint(tt.want) {
				t.Errorf("got %v, want %v", v, tt.want)
			}
		})
	}
}

func TestUriFragment_Access(t *testing.T) {
	var doc interface{}
	if err := json.Unmarshal([]byte(`{"foo":["bar","baz"],"":0,"a/b":1,"c%d":2,"e^f":3,"g|h":4,"i\\j":5,"k\"l":6," ":7,"m~n":8}`), &doc); err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		document interface{}
		uri      UriFragment
		want     interface{}
		err      error
	}{
		{doc, "#", doc, nil},
		{doc, "#/foo", []string{"bar", "baz"}, nil},
		{doc, "#/foo/0", "bar", nil},
		{doc, "#/", 0, nil},
		{doc, "#/a~1b", 1, nil},
		{doc, "#/c%25d", 2, nil},
		{doc, "#/e%5Ef", 3, nil},
		{doc, "#/g%7Ch", 4, nil},
		{doc, "#/i%5Cj", 5, nil},
		{doc, "#/k%22l", 6, nil},
		{doc, "#/%20", 7, nil},
		{doc, "#/m~0n", 8, nil},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v, err := tt.uri.Access(doc)
			if fmt.Sprint(err) != fmt.Sprint(tt.err) {
				t.Errorf("got error %v, want %v", err, tt.err)
			}
			if fmt.Sprint(v) != fmt.Sprint(tt.want) {
				t.Errorf("got %v, want %v", v, tt.want)
			}
		})
	}
}
