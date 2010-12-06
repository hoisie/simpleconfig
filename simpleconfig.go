package simpleconfig

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "reflect"
    "strconv"
    "strings"
)

func read(r io.Reader) (map[string]string, os.Error) {
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
        parts := strings.Split(line, " ", 2)
        if len(parts) == 2 {
            data[parts[0]] = parts[1]
        }
    }

    return data, nil
}

func Read(source interface{}) (map[string]string, os.Error) {
    if s, ok := source.(string); ok {
        //is the string a filename? if so, try opening it
        f, err := os.Open(s, os.O_RDONLY, 0666)
        if err == nil {
            //read config from the file
            defer f.Close()
            return read(f)
        } else {
            //read config from the string
            buffer := bytes.NewBufferString(s)
            return read(buffer)
        }
    }

    if r, ok := source.(io.Reader); ok {
        return read(r)
    }

    return nil, os.NewError("Invalid source type")
}


func writeTo(s string, val reflect.Value) os.Error {
    switch v := val.(type) {
    // if we're writing to an interace value, just set the byte data
    // TODO: should we support writing to a pointer?
    case *reflect.InterfaceValue:
        v.Set(reflect.NewValue(s))
    case *reflect.BoolValue:
        if strings.ToLower(s) == "false" || s == "0" {
            v.Set(false)
        } else {
            v.Set(true)
        }
    case *reflect.IntValue:
        i, err := strconv.Atoi64(s)
        if err != nil {
            return err
        }
        v.Set(i)
    case *reflect.UintValue:
        ui, err := strconv.Atoui64(s)
        if err != nil {
            return err
        }
        v.Set(ui)
    case *reflect.FloatValue:
        f, err := strconv.Atof64(s)
        if err != nil {
            return err
        }
        v.Set(f)

    case *reflect.StringValue:
        v.Set(s)
    case *reflect.SliceValue:
        typ := v.Type().(*reflect.SliceType)
        if _, ok := typ.Elem().(*reflect.UintType); ok {
            v.Set(reflect.NewValue([]byte(s)).(*reflect.SliceValue))
        }
    }
    return nil
}

// matchName returns true if key should be written to a field named name.
func matchName(key, name string) bool {
    return strings.ToLower(key) == strings.ToLower(name)
}

func writeToContainer(dst reflect.Value, data map[string]string) os.Error {
    switch v := dst.(type) {
    case *reflect.PtrValue:
        return writeToContainer(reflect.Indirect(v), data)
    case *reflect.InterfaceValue:
        return writeToContainer(v.Elem(), data)
    case *reflect.MapValue:
        if _, ok := v.Type().(*reflect.MapType).Key().(*reflect.StringType); !ok {
            return os.NewError("Invalid map key type")
        }
        elemtype := v.Type().(*reflect.MapType).Elem()
        for pk, pv := range data {
            mk := reflect.NewValue(pk)
            mv := reflect.MakeZero(elemtype)
            writeTo(pv, mv)
            v.SetElem(mk, mv)
        }
    case *reflect.StructValue:
        for pk, pv := range data {
            //try case sensitive match
            field := v.FieldByName(pk)
            if field != nil {
                writeTo(pv, field)
            }

            //try case insensitive matching
            field = v.FieldByNameFunc(func(s string) bool { return matchName(pk, s) })
            if field != nil {
                writeTo(pv, field)
            }

        }
    default:
        return os.NewError("Invalid container type")
    }
    return nil
}


func Unmarshal(dst interface{}, src interface{}) os.Error {
    data, err := Read(src)
    if err != nil {
        return err
    }

    return writeToContainer(reflect.NewValue(dst), data)
}
