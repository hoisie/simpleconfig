package simpleconfig

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func read(dst interface{}, r io.Reader) error {
	data := make(map[string]string)
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		//remove blank lines
		if len(line) == 0 {
			continue
		}
		//remove comment lines
		if strings.HasPrefix(line, "#") {
			continue
		}
		//config file is same with .properties file
		pos := strings.Index(line, "=")
		if pos > 0 {
			k := strings.Trim(line[:pos], " ")
			v := strings.Trim(line[pos+1:], " ")
			data[k] = v
		}
		//parts := strings.SplitN(line, " ", 2)
		//if len(parts) == 2 {
		//	data[parts[0]] = parts[1]
		//}
	}
	return unmarshal(dst, data)
}

func Read(dst interface{}, reader io.Reader) error {
	return read(dst, reader)
}

func ReadFile(dst interface{}, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	//read config from the file
	defer f.Close()
	return read(dst, f)
}

func ReadString(dst interface{}, s string) error {
	//read config from the string
	buffer := bytes.NewBufferString(s)
	return read(dst, buffer)
}

func writeTo(s string, val reflect.Value) error {
	switch v := val; v.Kind() {
	// if we're writing to an interace value, just set the byte data
	// TODO: should we support writing to a pointer?
	case reflect.Interface:
		v.Set(reflect.ValueOf(s))
	case reflect.Bool:
		if strings.ToLower(s) == "false" || s == "0" {
			v.SetBool(false)
		} else {
			v.SetBool(true)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(ui)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)

	case reflect.String:
		v.SetString(s)
	case reflect.Slice:
		typ := v.Type()
		if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
			v.Set(reflect.ValueOf([]byte(s)))
		}
	}
	return nil
}

// matchName returns true if key should be written to a field named name.
func matchName(key, name string) bool {
	return strings.ToLower(key) == strings.ToLower(name)
}

func writeToContainer(dst reflect.Value, data map[string]string) error {
	switch v := dst; v.Kind() {
	case reflect.Ptr:
		return writeToContainer(reflect.Indirect(v), data)
	case reflect.Interface:
		return writeToContainer(v.Elem(), data)
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return errors.New("Invalid map key type")
		}
		elemtype := v.Type().Elem()
		for pk, pv := range data {
			mk := reflect.ValueOf(pk)
			mv := reflect.Zero(elemtype)
			writeTo(pv, mv)
			v.SetMapIndex(mk, mv)
		}
	case reflect.Struct:
		for pk, pv := range data {
			//try case sensitive match
			field := v.FieldByName(pk)
			if field.IsValid() {
				writeTo(pv, field)
			}

			//try case insensitive matching
			field = v.FieldByNameFunc(func(s string) bool { return matchName(pk, s) })
			if field.IsValid() {
				writeTo(pv, field)
			}

		}
	default:
		return errors.New("Invalid container type")
	}
	return nil
}

func unmarshal(dst interface{}, data map[string]string) error {
	return writeToContainer(reflect.ValueOf(dst), data)
}
