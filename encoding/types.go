package encoding

// mess - still decodonging
const (
	// followed by byte string length (maybe varint, if string is >255?)
	typeString    = 0x30
	typeStringEnd = typeString // ?

	// followed by 4 bytes unix ts
	typeDate  = 0x40
	typeInt64 = 0x17

	// followed by byte n entries?
	typeMap = 0x60

	// followed by byte n entries
	typeArray = 0x50

	typeUnknownInt     = 0x11
	typeUnknownInt32   = 0x14 // probably uint32?
	typeUnknownInt32_2 = 0x15 // probably uint32?
	typeUnknown2Byte   = 0x13
	typeUint64         = 0x16

	typeUnknownNil       = 0x01 // ?
	typeUnknownBoolFalse = 0x03 //?
	typeUnknownBoolTrue  = 0x02 // ?

	typeUnknown1     = 0xC  // array ?
	typeUnknown5     = 0x10 // ?
	typeUnknown6     = 0x21 // ? some kind of hash, it's 8 bytes
	typeUnknown_bArr = 0x67 // ??? selected 0x52ad bytes, eh?
	// g4R. = 67 01 34 52 A9
	// 67 = byte[], 01 = len1?
	// 34 = type bytes, 52 A9 = len?
	typeUnknown_bytes = 0x34
)

type (
	AgValue interface{}

	AgMap    map[string]interface{}
	AgArray  map[string]interface{}
	AgString string
	AgDate   uint32
	AgInt64  int64
	AgInt    int
	AgUint32 uint32
	AgNil    int
	AgUint64 uint64

	AgUnknown         interface{}
	AgUnknown2Byte    []byte
	AgUnknown8Byte    []byte
	AgUnknown_ByteArr map[string]interface{}
	AgUnknown_Bytes   []byte
)
