package pkg

import "github.com/pkg/errors"

const (
	STREAM_MAGIC      uint16 = uint16(0xaced)
	STREAM_VERSION    uint16 = 5
	TC_NULL           byte   = 0x70 // 空指针
	TC_REFERENCE      byte   = 0x71
	TC_CLASSDESC      byte   = 0x72 // TC_CLASSDESC. 指定这是一个新类。
	TC_OBJECT         byte   = 0x73 // TC_OBJECT. 指定这是一个新的Object.
	TC_STRING         byte   = 0x74
	TC_ARRAY          byte   = 0x75
	TC_CLASS          byte   = 0x76
	TC_BLOCKDATA      byte   = 0x77
	TC_ENDBLOCKDATA   byte   = 0x78
	TC_RESET          byte   = 0x79
	TC_BLOCKDATALONG  byte   = 0x7A
	TC_EXCEPTION      byte   = 0x7B
	TC_LONGSTRING     byte   = 0x7C
	TC_PROXYCLASSDESC byte   = 0x7D
	TC_ENUM           byte   = 0x7E

	baseWireHandle    int  = 0x7E0000 /* First wire handle to be assigned. */
	SC_WRITE_METHOD   byte = 0x01     //if SC_SERIALIZABLE /* Flag bits for ObjectStreamClasses in Stream. */
	SC_BLOCK_DATA     byte = 0x08     //if SC_EXTERNALIZABLE
	SC_SERIALIZABLE   byte = 0x02
	SC_EXTERNALIZABLE byte = 0x04
	SC_ENUM           byte = 0x10
	SC_Fail           byte = 0x00
)

/*
  `B'       // byte
  `C'       // char
  `D'       // double
  `F'       // float
  `I'       // integer
  `J'       // long
  `S'       // short
  `Z'       // boolean
  `[`   	// array
  `L'       // object

*/
// primitiveHandlers maps serialized primitive identifiers to a corresponding primitiveHandler.
var primitiveHandlers = map[string]primitiveHandler{
	"B": func(sop *SerializedObjectParser) (b interface{}, err error) {
		if b, err = sop.readInt8(); err != nil {
			err = errors.Wrap(err, "error reading byte primitive")
		}

		return
	},
	"C": func(sop *SerializedObjectParser) (char interface{}, err error) {
		var charCode uint16
		if charCode, err = sop.readUInt16(); err != nil {
			err = errors.Wrap(err, "error reading char primitive")
		} else {
			char = string(rune(charCode))
		}

		return
	},
	"D": func(sop *SerializedObjectParser) (double interface{}, err error) {
		if double, err = sop.readFloat64(); err != nil {
			err = errors.Wrap(err, "error reading double primitive")
		}

		return
	},
	"F": func(sop *SerializedObjectParser) (f32 interface{}, err error) {
		if f32, err = sop.readFloat32(); err != nil {
			err = errors.Wrap(err, "error reading float primitive")
		}

		return
	},
	"I": func(sop *SerializedObjectParser) (i32 interface{}, err error) {
		if i32, err = sop.readInt32(); err != nil {
			err = errors.Wrap(err, "error reading int primitive")
		}

		return
	},
	"J": func(sop *SerializedObjectParser) (long interface{}, err error) {
		if long, err = sop.readInt64(); err != nil {
			err = errors.Wrap(err, "error reading long primitive")
		}

		return
	},
	"S": func(sop *SerializedObjectParser) (short interface{}, err error) {
		if short, err = sop.readInt16(); err != nil {
			err = errors.Wrap(err, "error reading short primitive")
		}

		return
	},
	"Z": func(sop *SerializedObjectParser) (b interface{}, err error) {
		var x int8
		if x, err = sop.readInt8(); err != nil {
			err = errors.Wrap(err, "error reading boolean primitive")
		} else {
			b = x != 0
		}

		return
	},
	"L": func(sop *SerializedObjectParser) (obj interface{}, err error) {
		if obj, err = sop.content(nil); err != nil {
			err = errors.Wrap(err, "error reading object primitive")
		}

		return
	},
	"[": func(sop *SerializedObjectParser) (arr interface{}, err error) {
		if arr, err = sop.content(nil); err != nil {
			err = errors.Wrap(err, "error reading array primitive")
		}

		return
	},
}
