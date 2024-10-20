// Package jsonpointer implements JSON Pointer according to RFC 6901
//
// See https://www.rfc-editor.org/rfc/rfc6901
package jsonpointer

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const ptrSep = "/"

var (
	jsonPointerRe = regexp.MustCompile(`^(/([^/~]|(~[01]))*)*$`)
	arrayIndexRe  = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)
)

var (
	escaper   = strings.NewReplacer("/", "~1", "~", "~0")
	unEscaper = strings.NewReplacer("~1", "/", "~0", "~")
)

type (
	Ptr         string
	UriFragment string
	token       struct {
		key string
		idx int
	}
)

type AccessError struct {
	used      []token
	remaining []token
	value     reflect.Value
	msg       string
}

func (e *AccessError) Error() string {
	return fmt.Sprintf("jsonpointer.Ptr: used %v, remaining %v, current val %v: %s", e.used, e.remaining, e.value, e.msg)
}

func (p Ptr) tokens() (tokens []token, err error) {
	s := string(p)
	if !jsonPointerRe.MatchString(s) {
		return nil, errors.New("jsonpointer.Ptr: invalid syntax")
	}

	parts := strings.Split(s, ptrSep)
	if len(parts) == 1 {
		return
	}
	for _, part := range parts[1:] {
		part = unEscaper.Replace(part)
		if idx, err := strconv.Atoi(part); err == nil {
			tokens = append(tokens, token{"", idx})
		} else {
			tokens = append(tokens, token{part, -1})
		}
	}
	return
}

func (p Ptr) Access(document any) (v reflect.Value, err error) {
	tokens, err := p.tokens()
	if err != nil {
		return reflect.Value{}, err
	}

	v = reflect.ValueOf(document)
	for i, t := range tokens {
		for {
			switch v.Kind() {
			case reflect.Array, reflect.Slice:
				if t.idx != -1 {
					if t.idx < v.Len() {
						v = v.Index(t.idx)
						goto next
					}
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "index out of range"}
				} else {
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "expect map"}
				}
			case reflect.Map:
				if t.idx == -1 {
					u := v.MapIndex(reflect.ValueOf(t.key))
					if u.IsValid() {
						v = u
						goto next
					}
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "key out of range"}
				} else {
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "expect array got object"}
				}
			case reflect.Struct:
				if t.idx == -1 {
					if i := fieldByJsonTag(v.Type(), t.key); i != -1 {
						v = v.Field(i)
						goto next
					}
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "key out of range"}
				} else {
					return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, "expect array got object"}
				}
			case reflect.Interface, reflect.Ptr:
				v = v.Elem()
			default:
				return reflect.Value{}, &AccessError{tokens[:i], tokens[i:], v, fmt.Sprintf("expect JSON got %v", v.Type())}
			}
		}
	next:
	}
	return v, nil
}

func fieldByJsonTag(t reflect.Type, key string) int {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name, _ := parseTag(tag)
		if !isValidTag(name) {
			name = f.Name
		}
		if name == key {
			return i
		}
	}
	return -1
}

func (uri UriFragment) Access(document any) (v reflect.Value, err error) {
	s := string(uri)
	if s, err = url.QueryUnescape(s); err != nil {
		return reflect.Value{}, err
	}
	if len(s) == 0 || s[0] != '#' {
		return reflect.Value{}, errors.New("jsonpointer.URIFragment: invalid syntax")
	}
	return Ptr(s[1:]).Access(document)
}
