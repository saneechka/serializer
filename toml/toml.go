package toml

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type TOMLSerializer struct{}

func New() *TOMLSerializer {
	return &TOMLSerializer{}
}

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenString
	tokenNumber
	tokenTrue
	tokenFalse
	tokenDate
	tokenLeftBracket
	tokenRightBracket
	tokenDot
	tokenEquals
	tokenComma
	tokenNewline
)

type token struct {
	typ   tokenType
	value string
}

type lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func newLexer(input string) *lexer {
	return &lexer{input: input, line: 1, col: 1}
}

func (l *lexer) next() token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return token{typ: tokenEOF}
	}

	switch c := l.input[l.pos]; c {
	case '[':
		l.pos++
		l.col++
		return token{typ: tokenLeftBracket, value: "["}
	case ']':
		l.pos++
		l.col++
		return token{typ: tokenRightBracket, value: "]"}
	case '.':
		l.pos++
		l.col++
		return token{typ: tokenDot, value: "."}
	case '=':
		l.pos++
		l.col++
		return token{typ: tokenEquals, value: "="}
	case ',':
		l.pos++
		l.col++
		return token{typ: tokenComma, value: ","}
	case '\n':
		l.pos++
		l.line++
		l.col = 1
		return token{typ: tokenNewline, value: "\n"}
	case '"':
		return l.readString()
	case 't':
		if l.pos+3 < len(l.input) && l.input[l.pos:l.pos+4] == "true" {
			l.pos += 4
			l.col += 4
			return token{typ: tokenTrue, value: "true"}
		}
	case 'f':
		if l.pos+4 < len(l.input) && l.input[l.pos:l.pos+5] == "false" {
			l.pos += 5
			l.col += 5
			return token{typ: tokenFalse, value: "false"}
		}
	}

	if c := l.input[l.pos]; c == '-' || unicode.IsDigit(rune(c)) {
		return l.readNumberOrDate()
	}

	return token{typ: tokenEOF}
}

func (l *lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if c == '\n' {
			l.line++
			l.col = 1
		} else if !unicode.IsSpace(rune(c)) {
			break
		}
		l.pos++
		l.col++
	}
}

func (l *lexer) readString() token {
	start := l.pos
	l.pos++ // skip opening quote
	l.col++

	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if c == '"' && l.input[l.pos-1] != '\\' {
			l.pos++ // skip closing quote
			l.col++
			return token{typ: tokenString, value: l.input[start+1 : l.pos-1]}
		}
		if c == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}

	return token{typ: tokenEOF}
}

func (l *lexer) readNumberOrDate() token {
	start := l.pos
	isDate := false

	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if c == 'T' || c == 'Z' || c == '-' || c == ':' {
			isDate = true
		} else if !unicode.IsDigit(rune(c)) && c != '.' && c != '+' && c != 'e' && c != 'E' {
			break
		}
		l.pos++
		l.col++
	}

	value := l.input[start:l.pos]
	if isDate {
		return token{typ: tokenDate, value: value}
	}
	return token{typ: tokenNumber, value: value}
}

type parser struct {
	lexer *lexer
	token token
}

func newParser(input string) *parser {
	lexer := newLexer(input)
	return &parser{
		lexer: lexer,
		token: lexer.next(),
	}
}

func (p *parser) next() {
	p.token = p.lexer.next()
}

func (p *parser) parseValue() (interface{}, error) {
	switch p.token.typ {
	case tokenString:
		val := p.token.value
		p.next()
		return val, nil
	case tokenNumber:
		val := p.token.value
		p.next()
		if strings.Contains(val, ".") {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			return f, nil
		}
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case tokenDate:
		val := p.token.value
		p.next()
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		return t, nil
	case tokenTrue:
		p.next()
		return true, nil
	case tokenFalse:
		p.next()
		return false, nil
	case tokenLeftBracket:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.token)
	}
}

func (p *parser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)
	p.next() // skip [

	if p.token.typ == tokenRightBracket {
		p.next()
		return arr, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		if p.token.typ == tokenRightBracket {
			p.next()
			return arr, nil
		}

		if p.token.typ != tokenComma {
			return nil, fmt.Errorf("expected comma or ], got %v", p.token)
		}
		p.next()
	}
}

func (p *parser) parseTable() (map[string]interface{}, error) {
	table := make(map[string]interface{})
	current := table
	path := make([]string, 0)

	for p.token.typ != tokenEOF {
		switch p.token.typ {
		case tokenLeftBracket:
			p.next()
			path = p.parseTablePath()
			if p.token.typ != tokenRightBracket {
				return nil, fmt.Errorf("expected ], got %v", p.token)
			}
			p.next()

			// Navigate to the correct nested map
			current = table
			for i, key := range path[:len(path)-1] {
				if _, exists := current[key]; !exists {
					current[key] = make(map[string]interface{})
				}
				if next, ok := current[key].(map[string]interface{}); ok {
					current = next
				} else {
					return nil, fmt.Errorf("cannot use %s as table, it's already defined as a value", strings.Join(path[:i+1], "."))
				}
			}
			current = current[path[len(path)-1]].(map[string]interface{})

		case tokenString:
			key := p.token.value
			p.next()

			if p.token.typ != tokenEquals {
				return nil, fmt.Errorf("expected =, got %v", p.token)
			}
			p.next()

			value, err := p.parseValue()
			if err != nil {
				return nil, err
			}
			current[key] = value

			if p.token.typ != tokenNewline && p.token.typ != tokenEOF {
				return nil, fmt.Errorf("expected newline or EOF, got %v", p.token)
			}
			p.next()

		default:
			return nil, fmt.Errorf("unexpected token: %v", p.token)
		}
	}

	return table, nil
}

func (p *parser) parseTablePath() []string {
	var path []string
	for {
		if p.token.typ != tokenString {
			return path
		}
		path = append(path, p.token.value)
		p.next()

		if p.token.typ != tokenDot {
			return path
		}
		p.next()
	}
}

func (s *TOMLSerializer) Unmarshal(data []byte, v any) error {
	parser := newParser(string(data))
	value, err := parser.parseTable()
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("v must be a pointer")
	}

	return s.setValue(rv.Elem(), value)
}

func (s *TOMLSerializer) setValue(rv reflect.Value, value interface{}) error {
	switch rv.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			rv.SetString(str)
		} else {
			return fmt.Errorf("cannot convert %v to string", value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case float64:
			rv.SetInt(int64(v))
		case int64:
			rv.SetInt(v)
		default:
			return fmt.Errorf("cannot convert %v to int", value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := value.(type) {
		case float64:
			rv.SetUint(uint64(v))
		case int64:
			rv.SetUint(uint64(v))
		default:
			return fmt.Errorf("cannot convert %v to uint", value)
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := value.(float64); ok {
			rv.SetFloat(f)
		} else {
			return fmt.Errorf("cannot convert %v to float", value)
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			rv.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %v to bool", value)
		}
	case reflect.Slice:
		if value == nil {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		arr, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("cannot convert %v to slice", value)
		}
		rv.Set(reflect.MakeSlice(rv.Type(), len(arr), len(arr)))
		for i, v := range arr {
			if err := s.setValue(rv.Index(i), v); err != nil {
				return err
			}
		}
	case reflect.Map:
		if value == nil {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		obj, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot convert %v to map", value)
		}
		rv.Set(reflect.MakeMap(rv.Type()))
		for k, v := range obj {
			key := reflect.ValueOf(k)
			elem := reflect.New(rv.Type().Elem()).Elem()
			if err := s.setValue(elem, v); err != nil {
				return err
			}
			rv.SetMapIndex(key, elem)
		}
	case reflect.Struct:
		if value == nil {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		obj, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot convert %v to struct", value)
		}
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			tomlTag := field.Tag.Get("toml")
			if tomlTag == "-" {
				continue
			}
			name := field.Name
			if tomlTag != "" {
				name = strings.Split(tomlTag, ",")[0]
			}
			if v, ok := obj[name]; ok {
				if err := s.setValue(rv.Field(i), v); err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		if value == nil {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return s.setValue(rv.Elem(), value)
	case reflect.Interface:
		if value == nil {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		rv.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("unsupported type: %v", rv.Kind())
	}
	return nil
}

func (s *TOMLSerializer) Marshal(v any) ([]byte, error) {
	return s.marshalValue(reflect.ValueOf(v), "")
}

func (s *TOMLSerializer) marshalValue(v reflect.Value, prefix string) ([]byte, error) {
	switch v.Kind() {
	case reflect.String:
		return []byte(`"` + escapeString(v.String()) + `"`), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(v.Bool())), nil
	case reflect.Slice, reflect.Array:
		return s.marshalArray(v)
	case reflect.Map:
		return s.marshalMap(v, prefix)
	case reflect.Struct:
		return s.marshalStruct(v, prefix)
	case reflect.Ptr:
		if v.IsNil() {
			return []byte("null"), nil
		}
		return s.marshalValue(v.Elem(), prefix)
	case reflect.Interface:
		if v.IsNil() {
			return []byte("null"), nil
		}
		return s.marshalValue(v.Elem(), prefix)
	case reflect.Invalid:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", v.Kind())
	}
}

func (s *TOMLSerializer) marshalArray(v reflect.Value) ([]byte, error) {
	if v.IsNil() {
		return []byte("[]"), nil
	}

	var elements []string
	for i := 0; i < v.Len(); i++ {
		element, err := s.marshalValue(v.Index(i), "")
		if err != nil {
			return nil, err
		}
		elements = append(elements, string(element))
	}
	return []byte("[" + strings.Join(elements, ", ") + "]"), nil
}

func (s *TOMLSerializer) marshalMap(v reflect.Value, prefix string) ([]byte, error) {
	if v.IsNil() {
		return []byte("{}"), nil
	}

	var pairs []string
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyStr := key.String()
		if prefix != "" {
			keyStr = prefix + "." + keyStr
		}

		valueBytes, err := s.marshalValue(value, keyStr)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, keyStr+" = "+string(valueBytes))
	}
	return []byte(strings.Join(pairs, "\n")), nil
}

func (s *TOMLSerializer) marshalStruct(v reflect.Value, prefix string) ([]byte, error) {
	var pairs []string
	t := v.Type()

	// Handle time.Time specially
	if t == reflect.TypeOf(time.Time{}) {
		t := v.Interface().(time.Time)
		return []byte(t.Format(time.RFC3339)), nil
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get TOML tag or use field name
		tomlTag := field.Tag.Get("toml")
		if tomlTag == "-" {
			continue
		}

		name := field.Name
		if tomlTag != "" {
			name = strings.Split(tomlTag, ",")[0]
		}

		if prefix != "" {
			name = prefix + "." + name
		}

		valueBytes, err := s.marshalValue(value, name)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, name+" = "+string(valueBytes))
	}
	return []byte(strings.Join(pairs, "\n")), nil
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

func (s *TOMLSerializer) Format() string {
	return "TOML"
}
