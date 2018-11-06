package encoding

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Decoder struct {
	r   *os.File
	buf []byte
}

func NewDecoder(r *os.File) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) DecodeRaw() (AgValue, error) {
	return d.getValue()
}

func (d *Decoder) Dump() {
	v, err := d.getValue()
	if err != nil {
		return
	}

	spew.Dump(v)
}

// i'm lazy ok
func (d *Decoder) pos() int {
	pos, _ := d.r.Seek(0, io.SeekCurrent)

	return int(pos)
}

func (d *Decoder) getValue() (AgValue, error) {
	buf := []byte{0}
	if _, err := d.r.Read(buf); err != nil {
		return nil, err
	}

	var v AgValue
	var token = buf[0]

	switch token {
	case typeMap:
		v = make(AgMap)
	case typeString:
		v = AgString("")
	case typeDate:
		v = AgDate(0)
	case typeInt64:
		v = AgInt64(0)
	case typeArray:
		v = make(AgArray)
	case typeUnknownInt32, typeUnknownInt32_2:
		v = AgUint32(0)
	case typeUint64:
		v = AgUint64(0)
	case typeUnknownInt:
		v = AgInt(0)
	case typeUnknown2Byte:
		v = AgUnknown2Byte(nil)
	case typeUnknown5:
		v = AgUnknown(0)

	case typeUnknown6:
		v = AgUnknown8Byte{}

	case typeUnknown_bArr:
		v = make(AgUnknown_ByteArr)
	case typeUnknown_bytes:
		v = AgUnknown_Bytes{}

	// none of these types have values after, so I assume part of the value
	// is the token; names are suspect, obviously, not definitive
	case typeUnknownNil:
		return nil, nil
	case typeUnknownBoolTrue:
		return true, nil
	case typeUnknownBoolFalse:
		return false, nil

	default:
		log.Printf("unknown type=%X at offset=%d (%x)", token, d.pos()-1, d.pos()-1)
	}

	v = d.readValue(v)

	return v, nil
}

func (d *Decoder) readValue(v AgValue) AgValue {
	switch v.(type) {
	case AgMap:
		v = d.readMap()
	case AgString:
		v = d.readString()
	case AgUint32:
		v = d.readUint32()
	case AgDate:
		v = d.readDate()
	case AgInt64:
		v = d.readInt64()
	case AgArray:
		v = d.readArray()
	case AgUint64:
		v = d.readUint64()
	case AgUnknown_ByteArr:
		v = d.readByteArr()
	case AgUnknown_Bytes:
		v = d.readBytes()
	case AgUnknown2Byte:
		v = d.readUnknown2Byte()
	case AgUnknown8Byte:
		v = d.readUnknown8Byte()
	case AgUnknown:
		v = d.readUnknown()
	case AgInt:
		v = d.readInt()
	case AgNil:
		v = nil
	default:
		//log.Printf("Don't know how to read value %q", v)
	}

	return v
}

func (d *Decoder) readString() AgString {
	var buf = []byte{0}

	d.r.Read(buf)
	n := int(buf[0])

	var str = make([]byte, n)
	d.r.Read(str)

	return AgString(string(str))
}

func (d *Decoder) readDate() AgDate {
	var buf = []byte{0, 0, 0, 0}

	d.r.Read(buf)

	return AgDate(binary.BigEndian.Uint32(buf))
}

func (d *Decoder) readInt64() AgInt64 {
	var buf = make([]byte, 8)
	d.r.Read(buf)

	return AgInt64(int64(binary.BigEndian.Uint64(buf)))
}

func (d *Decoder) readUint32() AgUint32 {
	var buf = make([]byte, 4)
	d.r.Read(buf)

	return AgUint32(binary.BigEndian.Uint32(buf))
}

func (d *Decoder) readUint64() AgUint64 {
	var buf = make([]byte, 8)
	d.r.Read(buf)

	return AgUint64(binary.BigEndian.Uint64(buf))
}

func (d *Decoder) readInt() AgInt {
	var buf = []byte{0}
	d.r.Read(buf)

	return AgInt(int(buf[0]))
}

func (d *Decoder) readUnknown() AgUnknown {
	var buf = []byte{0}
	d.r.Read(buf)

	return AgUnknown(buf)
}

func (d *Decoder) readUnknown2Byte() AgUnknown2Byte {
	var buf = []byte{0, 0}
	d.r.Read(buf)

	return AgUnknown2Byte(buf)
}

func (d *Decoder) readUnknown8Byte() AgUnknown8Byte {
	var buf = make([]byte, 8)
	d.r.Read(buf)

	return AgUnknown8Byte(buf)
}

func (d *Decoder) readArray() AgArray {
	var buf = []byte{0}
	d.r.Read(buf)

	n := int(buf[0])
	a := make(AgArray, n)

	for i := 0; i < n; i++ {
		value, err := d.getValue()
		if err != nil {
			panic(err)
		}

		a[strconv.Itoa(i)] = value
	}

	return a
}

func (d *Decoder) readByteArr() AgUnknown_ByteArr {
	var buf = []byte{0}
	d.r.Read(buf)
	n := int(buf[0])
	a := make(AgUnknown_ByteArr)

	for i := 0; i < n; i++ {
		val, err := d.getValue()
		if err != nil {
			panic(err)
		}

		a[strconv.Itoa(i)] = val
	}

	return a
}

func (d *Decoder) readBytes() AgUnknown_Bytes {
	var buf = []byte{0, 0}
	d.r.Read(buf)

	n := binary.BigEndian.Uint16(buf)
	b := make(AgUnknown_Bytes, int(n))
	d.r.Read(b)

	return b
}

func (d *Decoder) readMap() AgMap {
	var buf = []byte{0}
	d.r.Read(buf)

	n := int(buf[0])
	m := make(AgMap, n)

	// need to descend and decode..
	for i := 0; i < n; i++ {
		key, err := d.getValue()
		if err != nil {
			panic(err)
		}

		value, err := d.getValue()
		if err != nil {
			panic(err)
		}

		if _, ok := key.(AgString); !ok {
			log.Fatalf("key not a string as expected, what gives? offset %d (0x%X)\nm = %v", d.pos(), d.pos(), m)
		}

		skey := string(key.(AgString))

		m[skey] = value
	}

	return m
}
