package encoding

import (
	"log"
	"reflect"
	"strconv"
)

type Encoder struct {
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(v interface{}) AgValue {
	return e.encodeValue(reflect.ValueOf(v))
}

func (e *Encoder) encodeValue(v reflect.Value) AgValue {
	var agv AgValue

	k := v.Kind()
	switch k {
	case reflect.Interface:
		if v.IsNil() {
			return agv
		}
		agv = e.encodeValue(v.Elem())
	case reflect.Map:
		agv = e.encodeMap(v)
	case reflect.String:
		agv = AgString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		agv = AgInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		agv = AgInt(v.Int())
	case reflect.Array, reflect.Slice:
		log.Printf("got array value %v", v)
		agv = e.encodeArray(v)
	case reflect.Struct:
	default:
		log.Printf("no idea what to do with %v (%s)", v, k)
		//panic(fmt.Sprintf("can't encode unsupported value: rv=%s, rk=%s, rt=%s", rv, rk, rt))
	}

	return agv
}

func (e *Encoder) encodeMap(v reflect.Value) AgMap {
	keys := v.MapKeys()

	m := make(AgMap)
	for i := range keys {
		m[keys[i].String()] = e.encodeValue(v.MapIndex(keys[i]))
	}

	return m
}

func (e *Encoder) encodeArray(v reflect.Value) AgArray {
	a := make(AgArray)

	for i := 0; i < v.Len(); i++ {
		a[strconv.Itoa(i)] = e.encodeValue(v.Index(i))
	}

	return a
}
