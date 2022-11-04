package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	//_ "encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	//_ "strings"
	//_ "time"

	"github.com/pkg/errors"
)

// ParseSerializedObject parses a serialized java object.
func ParseSerializedObject(buf []byte) (content []interface{}, err error) {
	option := SetMaxDataBlockSize(len(buf))
	this := NewSerializedObjectParser(bytes.NewReader(buf), option)

	this.parseStream()
	return nil, nil
}

// ParseSerializedObject parses a serialized java object from stream.
func (this *SerializedObjectParser) ParseSerializedObject() (content []interface{}, err error) {
	if err = this.magic(); err != nil {
		return
	}

	if err = this.version(); err != nil {
		return
	}

	for !this.end() {
		var nxt interface{}

		if nxt, err = this.content(nil); err != nil {
			if errors.Cause(err).Error() == io.EOF.Error() {
				err = errors.New("premature end of input")
			}

			return
		}

		content = append(content, nxt)
	}

	return
}

// ParseSerializedObjectMinimal parses a serialized java object and returns the minimal object representation
// (i.e. without all the class info, etc...).
func ParseSerializedObjectMinimal(buf []byte) (content []interface{}, err error) {
	if content, err = ParseSerializedObject(buf); err == nil {
		content = jsonFriendlyArray(content)
	}

	return
}

// ParseSerializedObjectMinimal parses a serialized java object from stream
// and returns the minimal object representation (i.e. without all the class info, etc...).
func (this *SerializedObjectParser) ParseSerializedObjectMinimal() (content []interface{}, err error) {
	if content, err = this.ParseSerializedObject(); err == nil {
		content = jsonFriendlyArray(content)
	}

	return
}

// jsonFriendlyObject recursively filters / formats object fields to be as simple / JSON-like as possible.
func jsonFriendlyObject(obj interface{}) (jsonObj interface{}) {
	if m, isMap := obj.(map[string]interface{}); isMap {
		jsonMap := jsonFriendlyMap(m)
		jsonObj = jsonMap

		// if we have a single "value" key or a post-processed value just promote the value
		if mVal, mValExists := jsonMap["value"]; mValExists {
			_, mRawExists := jsonMap["@"]
			if mRawExists || len(jsonMap) == 1 {
				jsonObj = mVal
			}
		}

		return
	}

	if arr, isArray := obj.([]interface{}); isArray {
		jsonObj = jsonFriendlyArray(arr)

		return
	}

	// default for raw / primitive fields
	return obj
}

// jsonFriendlyArray recursively filters / formats a deserialized array.
func jsonFriendlyArray(arrayObj []interface{}) (jsonArray []interface{}) {
	jsonArray = make([]interface{}, len(arrayObj))
	for idx, arrayMember := range arrayObj {
		jsonArray[idx] = jsonFriendlyObject(arrayMember)
	}

	return
}

// jsonFriendlyMap recursively filters / formats a deserialized map.
func jsonFriendlyMap(mapObj map[string]interface{}) (jsonMap map[string]interface{}) {
	jsonMap = make(map[string]interface{})

	for k, v := range mapObj {
		// filter out `extends` keyword which just contains internal inheritance hierarchy
		if k == "extends" {
			continue
		}
		// filter out internal class definitions
		if _, isClazz := v.(*clazz); !isClazz {
			jsonMap[k] = jsonFriendlyObject(v)
		}
	}

	return
}

func init() {
	knownParsers = map[string]parser{
		"Enum":          parseEnum,
		"BlockDataLong": parseBlockDataLong,
		"BlockData":     parseBlockData,
		"EndBlockData":  parseEndBlockData,
		"ClassDesc":     parseClassDesc,
		"Class":         parseClass,
		"Array":         parseArray,
		"LongString":    parseLongString,
		"String":        parseString,
		"Null":          parseNull,
		"Object":        parseObject,
		"Reference":     parseReference,
	}
}

// typeNames includes all known type names.
var typeNames = []string{
	"Null",
	"Reference",
	"ClassDesc",
	"Object",
	"String",
	"Array",
	"Class",
	"BlockData",
	"EndBlockData",
	"Reset",
	"BlockDataLong",
	"Exception",
	"LongString",
	"ProxyClassDesc",
	"Enum",
}

// typeNameMax is used to ensure an encountered type is known.
var typeNameMax = uint8(len(typeNames) - 1)

// allowedClazzNames includes all allowed names when parsing a class descriptor.
var allowedClazzNames = map[string]bool{
	"ClassDesc":      true,
	"ProxyClassDesc": true,
	"Null":           true,
	"Reference":      true,
}

// parser is a func capable of reading a single serialized type.
type parser func(this *SerializedObjectParser) (interface{}, error)

// knownParsers maps serialized names to corresponding parser implementations.
var knownParsers map[string]parser

// PostProc handlers are used to format deserialized objects for easier consumption.
type PostProc func(map[string]interface{}, []interface{}) (map[string]interface{}, error)

// KnownPostProcs maps serialized object signatures to PostProc implementations.
var KnownPostProcs = map[string]PostProc{
	"java.util.ArrayList@7881d21d99c7619d":  listPostProc,
	"java.util.ArrayDeque@207cda2e240da08b": listPostProc,
	"java.util.Hashtable@13bb0f25214ae4b8":  mapPostProc,
	"java.util.HashMap@0507dac1c31660d1":    mapPostProc,
	"java.util.EnumMap@065d7df7be907ca1":    enumMapPostProc,
	"java.util.HashSet@ba44859596b8b734":    hashSetPostProc,
	"java.util.Date@686a81014b597419":       datePostProc,
}

// primitiveHandler are used to read primitive values.
type primitiveHandler func(this *SerializedObjectParser) (interface{}, error)

// SetMaxDataBlockSize set the maximum size of the parsed data block,
// by default it is equal to the value of the buffer size bufio.Reader or size of bytes.Reader.
func SetMaxDataBlockSize(maxSize int) Option {
	return func(this *SerializedObjectParser) {
		this.maxDataBlockSize = maxSize
	}
}

// NewSerializedObjectParser reads serialized java objects from stream.
func NewSerializedObjectParser(rd io.Reader, options ...Option) *SerializedObjectParser {
	buf := bufio.NewReaderSize(rd, bufferSize)
	sop := &SerializedObjectParser{
		rd:                     buf,
		maxDataBlockSize:       buf.Size(),
		_handleValue:           0x7e0000,
		_data:                  Smooth{data: []byte{}},
		_classDataDescriptions: []*ClassDataDesc{},
		so:                     &SerObject{},
	}
	sop._data._p = sop

	for _, option := range options {
		option(sop)
	}

	return sop
}

func (this *SerializedObjectParser) intToHex(i int) string {
	var b1 = make([]byte, 4)
	binary.BigEndian.PutUint32(b1, uint32(i))
	return fmt.Sprintf("%02x", b1[0]) +
		fmt.Sprintf(" %02x", b1[1]) +
		fmt.Sprintf(" %02x", b1[2]) +
		fmt.Sprintf(" %02x", b1[3])
	//return fmt.Sprintf("%s", hex.EncodeToString(b1))
	//return fmt.Sprintf("%02x", byte((i&0xff000000)>>24)) +
	//	fmt.Sprintf(" %02x", byte((i&0xff0000)>>16)) +
	//	fmt.Sprintf(" %02x", byte((i&0xff00)>>8)) +
	//	fmt.Sprintf(" %02x", byte(i&0xff))
}

func (this *SerializedObjectParser) parseStream() {
	var b1, b2 byte

	//The stream may begin with an RMI packet type byte, print it if so
	if b1 = this._data.peek(); b1 != STREAM_MAGIC1 {
		b1 = this._data.pop()
		switch b1 {
		case RMI_Call:
			this.print("RMI Call - 0x50")
			break
		case RMI_ReturnData:
			this.print("RMI ReturnData - 0x51")
			break
		case RMI_Ping:
			this.print("RMI Ping - 0x52")
			break
		case RMI_PingAck:
			this.print("RMI PingAck - 0x53")
			break
		case RMI_DgcAck:
			this.print("RMI DgcAck - 0x54")
			break
		default:
			this.print("Unknown RMI packet type - 0x" + this.byteToHex(b1))
			break
		}
	}

	//Magic number, print and validate
	b1 = this._data.pop()
	b2 = this._data.pop()
	this.print("STREAM_MAGIC - 0x" + this.byteToHex(b1) + " " + this.byteToHex(b2))
	if b1 != STREAM_MAGIC1 || b2 != STREAM_MAGIC2 {
		this.print("Invalid STREAM_MAGIC, should be 0xac ed")
		return
	}

	//Serialization version
	b1 = this._data.pop()
	b2 = this._data.pop()
	this.print("STREAM_VERSION - 0x" + this.byteToHex(b1) + " " + this.byteToHex(b2))
	if b1 != SC_Fail || b2 != STREAM_VERSION {
		this.print("Invalid STREAM_VERSION, should be 0x00 05")
	}

	//Remainder of the stream consists of one or more 'content' elements
	this.print("Contents")
	this.increaseIndent()
	for this._data.size() > 0 {
		if nil != this.readContentElement() {
			break
		}
	}
	this.decreaseIndent()
}

func (this *SerializedObjectParser) newHandle1() int {
	handleValue := this._handleValue
	//Print the handle value
	this.print("newHandle 0x" + this.intToHex(handleValue))

	//Increment the next handle value and return the one we just assigned
	this._handleValue++
	return int(handleValue)
}

// newHandle adds a parsed object to the existing indexed handles which can be used later to lookup references to
// existing objects.
func (this *SerializedObjectParser) newHandle(obj interface{}) interface{} {
	this.handles = append(this.handles, obj)

	return obj
}

func (this *SerializedObjectParser) print(s ...interface{}) {
	fmt.Printf(this._indent)
	for _, x := range s {
		fmt.Printf("%v", x)
	}
	fmt.Println("")
}
func (this *SerializedObjectParser) byteToHex(s uint8) string {
	var data = []byte{s}
	return hex.EncodeToString(data)
}

func (this *SerializedObjectParser) increaseIndent() {
	this._indent = this._indent + "  "
}

func (this *SerializedObjectParser) readNewEnum() {
	var b1 uint8

	//TC_ENUM
	b1 = this._data.pop()
	this.print("TC_ENUM - 0x" + this.byteToHex(b1))
	if b1 != TC_ENUM {
		log.Panicln("Error: Illegal value for TC_ENUM (should be 0x7e)")
	}

	// Indent
	this.increaseIndent()

	// classDesc
	this.readClassDesc()

	// newHandle
	this.newHandle1()

	//enumConstantName
	this.readNewString()

	//Decrease indent
	this.decreaseIndent()
}

/*******************
 * Decrease the indentation string or trigger an exception if the length is
 * already below 2.
 ******************/
func (this *SerializedObjectParser) decreaseIndent() {
	if len(this._indent) < 2 {
		log.Panicln("Error: Illegal indentation decrease.")
	}
	this._indent = this._indent[0 : len(this._indent)-2]
}

func (this *SerializedObjectParser) readUtf() string {
	var content = ""
	var hex = ""
	var b1, b2 uint8
	var len int
	//length
	b1 = this._data.pop()
	b2 = this._data.pop()

	len = (int(b1<<8) & 0xff00) + int(b2&0xff)
	this.print("Length - ", len, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2))

	//Contents
	for i := 0; i < len; {
		i += 1
		b1 = this._data.pop()
		content += fmt.Sprintf("%c", b1)
		hex += this.byteToHex(b1)
	}
	this.print("Value - " + content + " - 0x" + hex)
	//Return the string
	return content
}
func (this *SerializedObjectParser) readTC_CLASSDESC() *ClassDataDesc {
	var cdd = NewClassDataDesc()
	var b1 uint8

	//TC_CLASSDESC
	b1 = this._data.pop()
	this.print("TC_CLASSDESC - 0x" + this.byteToHex(b1))
	if b1 != TC_CLASSDESC {
		log.Panicln("Error: Illegal value for TC_CLASSDESC (should be 0x72)")
	}
	this.increaseIndent()

	//className
	this.print("className")
	this.increaseIndent()
	cdd.addClass(this.readUtf()) //Add the class name to the class data description
	this.decreaseIndent()

	//serialVersionUID
	this.print("serialVersionUID - 0x" + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()) +
		" " + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()) + " " + this.byteToHex(this._data.pop()))

	//newHandle
	cdd.setLastClassHandle(this.newHandle1()) //Set the reference handle for the most recently added class

	//classDescInfo
	this.readClassDescInfo(cdd) //Read class desc info, add the super class description to the ClassDataDesc if one is found

	//Decrease the indent
	this.decreaseIndent()

	//Return the ClassDataDesc
	return cdd
}
func (this *SerializedObjectParser) readClassDescInfo(cdd *ClassDataDesc) {
	var classDescFlags = ""
	var b1 byte

	//classDescFlags
	b1 = this._data.pop()
	if (b1 & SC_WRITE_METHOD) == SC_WRITE_METHOD {
		classDescFlags += "SC_WRITE_METHOD | "
	}
	if (b1 & SC_SERIALIZABLE) == SC_SERIALIZABLE {
		classDescFlags += "SC_SERIALIZABLE | "
	}
	if (b1 & SC_EXTERNALIZABLE) == SC_EXTERNALIZABLE {
		classDescFlags += "SC_EXTERNALIZABLE | "
	}
	if (b1 & SC_BLOCK_DATA) == SC_BLOCK_DATA {
		classDescFlags += "SC_BLOCK_DATA | "
	}
	if len(classDescFlags) > 0 {
		classDescFlags = classDescFlags[0 : len(classDescFlags)-3]
	}
	this.print("classDescFlags - 0x" + this.byteToHex(b1) + " - " + classDescFlags)

	//Store the classDescFlags
	cdd.setLastClassDescFlags(b1) //Set the classDescFlags for the most recently added class

	//Validate classDescFlags
	if (b1 & SC_SERIALIZABLE) == SC_SERIALIZABLE {
		if (b1 & SC_EXTERNALIZABLE) == SC_EXTERNALIZABLE {
			log.Panicln("Error: Illegal classDescFlags, SC_SERIALIZABLE is not compatible with SC_EXTERNALIZABLE.")
		}
		if (b1 & SC_BLOCK_DATA) == SC_BLOCK_DATA {
			log.Panicln("Error: Illegal classDescFlags, SC_SERIALIZABLE is not compatible with SC_BLOCK_DATA.")
		}
	} else if (b1 & SC_EXTERNALIZABLE) == SC_EXTERNALIZABLE {
		if (b1 & SC_WRITE_METHOD) == SC_WRITE_METHOD {
			log.Panicln("Error: Illegal classDescFlags, SC_EXTERNALIZABLE is not compatible with SC_WRITE_METHOD.")
		}
	} else if b1 != SC_Fail {
		log.Panicln("Error: Illegal classDescFlags, must include either SC_SERIALIZABLE or SC_EXTERNALIZABLE.")
	}

	//fields
	this.readFields(cdd) //Read field descriptions and add them to the ClassDataDesc

	//classAnnotation
	this.readClassAnnotation()

	//superClassDesc
	cdd.addSuperClassDesc(this.readSuperClassDesc()) //Read the super class description and add it to the ClassDataDesc
}

/*******************
 * Read a superClassDesc from the stream.
 *
 * classDesc
 ******************/
func (this *SerializedObjectParser) readSuperClassDesc() *ClassDataDesc {
	var cdd *ClassDataDesc
	//Print header, indent, delegate, decrease indent
	this.print("superClassDesc")
	this.increaseIndent()
	cdd = this.readClassDesc()
	this.decreaseIndent()

	//Return the super class data description
	return cdd
}

/*******************
 * Read a classAnnotation from the stream.
 *
 * Could be either:
 *	contents
 *	TC_ENDBLOCKDATA		(0x78)
 ******************/
func (this *SerializedObjectParser) readClassAnnotation() {
	//Print annotation section and indent
	this.print("classAnnotations")
	this.increaseIndent()

	//Loop until we have a TC_ENDBLOCKDATA
	x := this._data.peek()
	for x != TC_ENDBLOCKDATA {
		// Read a content element
		this.readContentElement()
		x = this._data.peek()
	}

	//Pop and print the TC_ENDBLOCKDATA element
	this._data.pop()
	this.print("TC_ENDBLOCKDATA - 0x78")

	//Decrease indent
	this.decreaseIndent()
}

/*******************
 * Read a content element from the data stream.
 *
 * Could be any of:
 *	TC_OBJECT			(0x73)
 *	TC_CLASS			(0x76)
 *	TC_ARRAY			(0x75)
 *	TC_STRING			(0x74)
 *	TC_LONGSTRING		(0x7c)
 *	TC_ENUM				(0x7e)
 *	TC_CLASSDESC		(0x72)
 *	TC_PROXYCLASSDESC	(0x7d)
 *	TC_REFERENCE		(0x71)
 *	TC_NULL				(0x70)
 *	TC_EXCEPTION		(0x7b)
 *	TC_RESET			(0x79)
 *	TC_BLOCKDATA		(0x77)
 *	TC_BLOCKDATALONG	(0x7a)
 ******************/
func (this *SerializedObjectParser) readContentElement() error {
	//Peek the next byte and delegate to the appropriate method
	switch this._data.peek() {
	case TC_OBJECT: //TC_OBJECT
		this.readNewObject()
		break

	case TC_CLASS: //TC_CLASS
		this.readNewClass()
		break

	case TC_ARRAY: //TC_ARRAY
		this.readNewArray()
		break

	case TC_STRING:
	//TC_STRING
	case TC_LONGSTRING: //TC_LONGSTRING
		this.readNewString()
		break

	case TC_ENUM: //TC_ENUM
		this.readNewEnum()
		break

	case TC_CLASSDESC:
	//TC_CLASSDESC
	case TC_PROXYCLASSDESC: //TC_PROXYCLASSDESC
		this.readNewClassDesc()
		break

	case TC_REFERENCE: //TC_REFERENCE
		this.readPrevObject()
		break

	case TC_NULL: //TC_NULL
		this.readNullReference()
		break

	case TC_EXCEPTION: //TC_EXCEPTION
		this.readException()
		break

	case TC_RESET: //TC_RESET
		this.handleReset()
		break

	case TC_BLOCKDATA: //TC_BLOCKDATA
		this.readBlockData()
		break

	case TC_BLOCKDATALONG: //TC_BLOCKDATALONG
		this.readLongBlockData()
		break

	default:
		//this.print("Invalid content element type 0x" + this.byteToHex(this._data.peek()))
		return errors.New("Error: Illegal content element type.")
	}
	return nil
}

/*******************
 * Read a fieldDesc from the stream.
 *
 * Could be either:
 *	prim_typecode	fieldName
 *	obj_typecode	fieldName	className1
 ******************/
func (this *SerializedObjectParser) readFieldDesc(cdd *ClassDataDesc) {
	var b1 byte
	//prim_typecode/obj_typecode
	b1 = this._data.pop()
	cdd.addFieldToLastClass(b1) //Add a field of the type in b1 to the most recently added class
	switch b1 {
	case 'B': //byte
		this.print("Byte - B - 0x" + this.byteToHex(b1))
		break

	case 'C': //char
		this.print("Char - C - 0x" + this.byteToHex(b1))
		break

	case 'D': //double
		this.print("Double - D - 0x" + this.byteToHex(b1))
		break

	case 'F': //float
		this.print("Float - F - 0x" + this.byteToHex(b1))
		break

	case 'I': //int
		this.print("Int - I - 0x" + this.byteToHex(b1))
		break

	case 'J': //long
		this.print("Long - L - 0x" + this.byteToHex(b1))
		break

	case 'S': //Short
		this.print("Short - S - 0x" + this.byteToHex(b1))
		break

	case 'Z': //boolean
		this.print("Boolean - Z - 0x" + this.byteToHex(b1))
		break

	case '[': //array
		this.print("Array - [ - 0x" + this.byteToHex(b1))
		break

	case 'L': //object
		this.print("Object - L - 0x" + this.byteToHex(b1))
		break

	default:
		//Unknown field type code
		log.Panicf("Error: Illegal field type code ('%c', 0x"+this.byteToHex(b1)+")", b1)
	}

	//fieldName
	this.print("fieldName")
	this.increaseIndent()
	cdd.setLastFieldName(this.readUtf()) //Set the name of the most recently added field
	this.decreaseIndent()

	//className1 (if non-primitive type)
	if b1 == '[' || b1 == 'L' {
		this.print("className1")
		this.increaseIndent()
		cdd.setLastFieldClassName1(this.readNewString()) //Set the className1 of the most recently added field
		this.decreaseIndent()
	}
}

/*******************
 * Read a newString element from the stream.
 *
 * Could be:
 *	TC_STRING		(0x74)
 *	TC_LONGSTRING	(0x7c)
 ******************/
func (this *SerializedObjectParser) readNewString() string {

	//Peek the type and delegate to the appropriate method
	switch this._data.peek() {
	case TC_STRING: //TC_STRING
		return this.readTC_STRING()

	case TC_LONGSTRING: //TC_LONGSTRING
		return this.readTC_LONGSTRING()

	case TC_REFERENCE: //TC_REFERENCE
		this.readPrevObject()
		return "[TC_REF]"

	default:
		this.print("Invalid newString type 0x" + this.byteToHex(this._data.peek()))
		log.Panicf("Error illegal newString type.")
	}
	return ""
}

/*******************
 * Read a TC_LONGSTRING element from the stream.
 *
 * TC_LONGSTRING	newHandle	long-utf
 ******************/
func (this *SerializedObjectParser) readTC_LONGSTRING() string {
	var val string
	var b1 byte

	//TC_LONGSTRING
	b1 = this._data.pop()
	this.print("TC_LONGSTRING - 0x" + this.byteToHex(b1))
	if b1 != TC_LONGSTRING {
		log.Panicln("Error: Illegal value for TC_LONGSTRING (should be 0x7c)")
	}

	//Indent
	this.increaseIndent()

	//newHandle
	this.newHandle1()

	//long-utf
	val = this.readLongUtf()

	//Decrease indent
	this.decreaseIndent()

	//Return the string value
	return val
}

/*******************
 * Read a long-UTF string from the stream.
 *
 * (long)length		contents
 ******************/
func (this *SerializedObjectParser) readLongUtf() string {
	var content = ""
	var hex = ""
	var b1, b2, b3, b4, b5, b6, b7, b8 byte
	var len uint64

	//Length
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	b5 = this._data.pop()
	b6 = this._data.pop()
	b7 = this._data.pop()
	b8 = this._data.pop()
	len = (uint64(b1<<56) & 0xff00000000000000) +
		(uint64(b2<<48) & 0xff000000000000) +
		(uint64(b3<<40) & 0xff0000000000) +
		(uint64(b4<<32) & 0xff00000000) +
		(uint64(b5<<24) & 0xff000000) +
		(uint64(b6<<16) & 0xff0000) +
		(uint64(b7<<8) & 0xff00) +
		uint64(b8&0xff)
	this.print("Length - ", len, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4)+" "+
		this.byteToHex(b5)+" "+this.byteToHex(b6)+" "+this.byteToHex(b7)+" "+this.byteToHex(b8))

	//Contents
	var l uint64 = 0
	for l < len {
		l += 1
		b1 = this._data.pop()
		content += fmt.Sprintf("%c", b1)
		hex += this.byteToHex(b1)
	}
	this.print("Value - " + content + " - 0x" + hex)

	//Return the string
	return content
}

func (this *SerializedObjectParser) readFields(cdd *ClassDataDesc) {
	var b1, b2 byte
	var count uint

	//count
	b1 = this._data.pop()
	b2 = this._data.pop()
	count = (uint(b1<<8) & 0xff00) + uint(b2&0xff)
	this.print("fieldCount - ", count, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2))

	//fieldDesc
	if count > 0 {
		this.print("Fields")
		this.increaseIndent()
		var i uint = 0
		for i < count {
			this.print(i, ":")
			this.increaseIndent()
			this.readFieldDesc(cdd)
			this.decreaseIndent()
			i += 1
		}
		this.decreaseIndent()
	}
}

/*******************
 * Read a newClassDesc from the stream.
 *
 * Could be:
 *	TC_CLASSDESC		(0x72)
 *	TC_PROXYCLASSDESC	(0x7d)
 ******************/
func (this *SerializedObjectParser) readNewClassDesc() *ClassDataDesc {
	var cdd *ClassDataDesc

	//Peek the type and delegate to the appropriate method
	switch this._data.peek() {
	case TC_CLASSDESC: //TC_CLASSDESC
		cdd = this.readTC_CLASSDESC()
		this._classDataDescriptions = append(this._classDataDescriptions, cdd)
		return cdd

	case TC_PROXYCLASSDESC: //TC_PROXYCLASSDESC
		cdd = this.readTC_PROXYCLASSDESC()
		this._classDataDescriptions = append(this._classDataDescriptions, cdd)
		return cdd

	default:
		this.print("Invalid newClassDesc type 0x" + this.byteToHex(this._data.peek()))
		log.Panicln("Error illegal newClassDesc type.")
	}
	return cdd
}

/*******************
 * Read a TC_PROXYCLASSDESC from the stream.
 *
 * TC_PROXYCLASSDESC	newHandle	proxyClassDescInfo
 ******************/
func (this *SerializedObjectParser) readTC_PROXYCLASSDESC() *ClassDataDesc {
	var cdd = NewClassDataDesc()
	var b1 byte

	//TC_PROXYCLASSDESC
	b1 = this._data.pop()
	this.print("TC_PROXYCLASSDESC - 0x" + this.byteToHex(b1))
	if b1 != TC_PROXYCLASSDESC {
		log.Panicln("Error: Illegal value for TC_PROXYCLASSDESC (should be 0x7d)")
	}
	this.increaseIndent()

	//Create the new class descriptor
	cdd.addClass("<Dynamic Proxy Class>")

	//newHandle
	cdd.setLastClassHandle(this.newHandle1()) //Set the reference handle for the most recently added class

	//proxyClassDescInfo
	this.readProxyClassDescInfo(cdd) //Read proxy class desc info, add the super class description to the ClassDataDesc if one is found

	//Decrease the indent
	this.decreaseIndent()

	//Return the ClassDataDesc
	return cdd
}

/*******************
 * Read a proxyClassDescInfo from the stream.
 *
 * (int)count	(utf)proxyInterfaceName[count]	classAnnotation		superClassDesc
 ******************/
func (this *SerializedObjectParser) readProxyClassDescInfo(cdd *ClassDataDesc) {
	var b1, b2, b3, b4 byte
	var count int

	//count
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	count = (int(b1<<24) & 0xff000000) +
		(int(b2<<16) & 0xff0000) +
		(int(b3<<8) & 0xff00) +
		(int(b4) & 0xff)
	this.print("Interface count - ", count, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4))

	//proxyInterfaceName[count]
	this.print("proxyInterfaceNames")
	this.increaseIndent()
	for i := 0; i < count; {
		i += 1
		this.print(i, ":")
		this.increaseIndent()
		this.readUtf()
		this.decreaseIndent()
	}
	this.decreaseIndent()

	//classAnnotation
	this.readClassAnnotation()

	//superClassDesc
	cdd.addSuperClassDesc(this.readSuperClassDesc()) //Read the super class description and add it to the ClassDataDesc
}

/*******************
 * Read a classDesc from the data stream.
 *
 * Could be:
 *	TC_CLASSDESC		(0x72)
 *	TC_PROXYCLASSDESC	(0x7d)
 *	TC_NULL				(0x70)
 *	TC_REFERENCE		(0x71)
 ******************/
func (this *SerializedObjectParser) readClassDesc() *ClassDataDesc {
	var refHandle int
	var b1 = this._data.peek()
	// Peek the type and delegate to the appropriate method
	switch b1 {
	// TC_CLASSDESC
	case TC_CLASSDESC, TC_PROXYCLASSDESC: // TC_PROXYCLASSDESC
		return this.readNewClassDesc()

	case TC_NULL: //TC_NULL
		this.readNullReference()
		return nil

	case TC_REFERENCE: //TC_REFERENCE
		refHandle = this.readPrevObject()                 //Look up a referenced class data description object and return it
		for _, cdd := range this._classDataDescriptions { //Iterate over all class data descriptions
			for classIndex := 0; classIndex < cdd.getClassCount(); { //Iterate over all classes in this class data description
				classIndex += 1
				if cdd.getClassDetails(classIndex).getHandle() == refHandle { //Check if the reference handle matches
					return cdd.buildClassDataDescFromIndex(classIndex) //Generate a ClassDataDesc starting from the given index and return it
				}
			}
		}
		//Invalid classDesc reference handle
		log.Panicln("Error: Invalid classDesc reference (0x" + this.intToHex(refHandle) + ")")

	default:
		this.print("Invalid classDesc type 0x" + this.byteToHex(this._data.peek()))
		log.Panicln("Error illegal classDesc type.")
	}
	return nil
}

// 读新对象
func (this *SerializedObjectParser) readNewObject() {
	var cdd *ClassDataDesc //ClassDataDesc describing the format of the objects 'classdata' element
	//TC_OBJECT
	b1 := this._data.pop()
	this.print("TC_OBJECT - 0x" + this.byteToHex(b1))
	if b1 != TC_OBJECT {
		log.Panicln("Error: Illegal value for TC_OBJECT (should be 0x73)")
	}

	// Indent
	this.increaseIndent()

	// classDesc
	cdd = this.readClassDesc() //Read the class data description

	//newHandle
	this.newHandle1()

	//classdata
	this.readClassData(cdd) //Read the class data based on the class data description - TODO This needs to check if cdd is null before reading anything

	//Decrease indent
	this.decreaseIndent()
}

func (this *SerializedObjectParser) readClassData(cdd *ClassDataDesc) {
	var cd *ClassDetails
	var classIndex int

	//Print header and indent
	this.print("classdata")
	this.increaseIndent()

	//Print class data if there is any
	if cdd != nil {
		//Iterate backwards through the classes as we need to deal with the most super (last added) class first
		for classIndex = cdd.getClassCount() - 1; classIndex >= 0; classIndex -= 1 {
			//Get the class details
			cd = cdd.getClassDetails(classIndex)

			//Print the class name and indent
			this.print(cd.getClassName())
			this.increaseIndent()

			//Read the field values if the class is SC_SERIALIZABLE
			if cd.isSC_SERIALIZABLE() {
				//Start the field values section and indent
				this.print("values")
				this.increaseIndent()

				//Iterate over the field and read/print the contents
				for _, cf := range cd.getFields() {
					this.readClassDataField(cf)
				}

				//Revert indent
				this.decreaseIndent()
			}

			hasBlockData := cd.isSC_SERIALIZABLE() && cd.isSC_WRITE_METHOD()

			if cd.isSC_EXTERNALIZABLE() {
				if cd.isSC_BLOCK_DATA() {
					hasBlockData = true
				} else { //Protocol version 1 does not use block data; cannot parse it
					this.increaseIndent()
					this.print("Unable to parse externalContents for protocol version 1.")
					log.Panicln("Error: Unable to parse externalContents element.")
				}
			}

			//Read object annotations
			if hasBlockData {
				//Start the object annotations section and indent
				this.print("objectAnnotation")
				this.increaseIndent()

				//Loop until we have a TC_ENDBLOCKDATA
				var x1 = this._data.peek()
				for x1 != TC_ENDBLOCKDATA {
					//Read a content element
					this.readContentElement()
					x1 = this._data.peek()
				}

				//Pop and print the TC_ENDBLOCKDATA element
				this._data.pop()
				this.print("TC_ENDBLOCKDATA - 0x78")

				//Revert indent
				this.decreaseIndent()
			}

			//Revert indent for this class
			this.decreaseIndent()
		}
	} else {
		this.print("N/A")
	}

	//Revert indent
	this.decreaseIndent()
}

/*******************
 * Read a classdata field from the stream.
 *
 * The data type depends on the given field description.
 *
 * @param f A description of the field data to read.
 ******************/
func (this *SerializedObjectParser) readClassDataField(cf *ClassField) {
	//Print the field name and indent
	this.print(cf.getName())
	this.increaseIndent()

	//Read the field data
	this.readFieldValue(cf.getTypeCode())

	//Decrease the indent
	this.decreaseIndent()
}

/*******************
 * Read a byte field.
 ******************/
func (this *SerializedObjectParser) readByteField() {
	var b1 byte = this._data.pop()
	c1 := fmt.Sprintf("%c", b1)
	if b1 >= 0x20 && b1 <= TC_ENUM {
		//Print with ASCII
		this.print("(byte)", c1, " (ASCII: "+c1+") - 0x"+this.byteToHex(b1))
	} else {
		//Just print byte value
		this.print("(byte)", b1, " - 0x"+this.byteToHex(b1))
	}
}

/*******************
 * Read a char field.
 ******************/
func (this *SerializedObjectParser) readCharField() {
	var b1 byte = this._data.pop()
	var b2 byte = this._data.pop()
	c1 := fmt.Sprintf("%c", byte((uint32(b1<<8)&0xff00)+uint32(b2&0xff)))
	this.print("(char)" + c1 + " - 0x" + this.byteToHex(b1) + " " + this.byteToHex(b2))
}

/*******************
 * Read a float field.
 ******************/
func (this *SerializedObjectParser) readFloatField() {
	var b1, b2, b3, b4 byte
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	var xx1 float64 = float64((uint32(b1<<24) & 0xff000000) +
		(uint32(b2<<16) & 0xff0000) +
		(uint32(b3<<8) & 0xff00) +
		uint32(b4&0xff))
	this.print("(float)", xx1, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+
		" "+this.byteToHex(b4))
}

/*******************
 * Read a field value based on the type code.
 *
 * @param typeCode The field type code.
 ******************/
func (this *SerializedObjectParser) readFieldValue(typeCode byte) {
	switch typeCode {
	case 'B': //byte
		this.readByteField()
		break

	case 'C': //char
		this.readCharField()
		break

	case 'D': //double
		this.readDoubleField()
		break

	case 'F': //float
		this.readFloatField()
		break

	case 'I': //int
		this.readIntField()
		break

	case 'J': //long
		this.readLongField()
		break

	case 'S': //short
		this.readShortField()
		break

	case 'Z': //boolean
		this.readBooleanField()
		break

	case '[': //array
		this.readArrayField()
		break

	case 'L': //object
		this.readObjectField()
		break

	default: //Unknown field type
		log.Panicln("Error: Illegal field type code ('", typeCode, "', 0x"+this.byteToHex(typeCode)+")")
	}
}

/*******************
 * Read an int field.
 ******************/
func (this *SerializedObjectParser) readIntField() {
	var b1, b2, b3, b4 byte
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	this.print("(int)", (int)((uint32(b1<<24)&0xff000000)+
		(uint32(b2<<16)&0xff0000)+
		(uint32(b3<<8)&0xff00)+
		uint32(b4&0xff)), " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+
		" "+this.byteToHex(b4))
}

/*******************
 * Read a long field.
 ******************/
func (this *SerializedObjectParser) readLongField() {
	var b1, b2, b3, b4, b5, b6, b7, b8 byte
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	b5 = this._data.pop()
	b6 = this._data.pop()
	b7 = this._data.pop()
	b8 = this._data.pop()
	this.print("(long)", (uint64(b1<<56)&0xff00000000000000)+
		(uint64(b2<<48)&0xff000000000000)+
		(uint64(b3<<40)&0xff0000000000)+
		(uint64(b4<<32)&0xff00000000)+
		(uint64(b5<<24)&0xff000000)+
		(uint64(b6<<16)&0xff0000)+
		(uint64(b7<<8)&0xff00)+
		uint64(b8&0xff), " - 0x"+this.byteToHex(b1)+
		" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4)+" "+this.byteToHex(b5)+" "+this.byteToHex(b6)+" "+
		this.byteToHex(b7)+" "+this.byteToHex(b8))
}

/*******************
 * Read a short field.
 ******************/
func (this *SerializedObjectParser) readShortField() {
	var b1, b2 byte
	b1 = this._data.pop()
	b2 = this._data.pop()
	this.print("(short)", uint16((uint16(b1<<8)&0xff00)+uint16(b2&0xff)), " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2))
}

/*******************
 * Read a boolean field.
 ******************/
func (this *SerializedObjectParser) readBooleanField() {
	var b1 = this._data.pop()
	var x1 = "true"
	if b1 == 0 {
		x1 = "false"
	}
	this.print("(boolean)" + x1 + " - 0x" + this.byteToHex(b1))
}

/*******************
 * Read an array field.
 ******************/
func (this *SerializedObjectParser) readArrayField() {
	//Print field type and increase indent
	this.print("(array)")
	this.increaseIndent()

	//Array could be null
	switch this._data.peek() {
	case TC_NULL: //
		this.readNullReference()
		break

	case TC_ARRAY: //TC_ARRAY
		this.readNewArray()
		break

	case TC_REFERENCE: //TC_REFERENCE
		this.readPrevObject()
		break

	default: //Unknown
		log.Panicln("Error: Unexpected array field value type (0x" + this.byteToHex(this._data.peek()))
	}

	//Revert indent
	this.decreaseIndent()
}

/*******************
 * Read an object field.
 ******************/
func (this *SerializedObjectParser) readObjectField() {
	this.print("(object)")
	this.increaseIndent()

	//Object fields can have various types of values...
	switch this._data.peek() {
	case TC_OBJECT: // TC_OBJECT New object
		this.readNewObject()
		break

	case TC_REFERENCE: // TC_REFERENCE
		this.readPrevObject()
		break

	case TC_NULL: // TC_NULL
		this.readNullReference()
		break

	case TC_STRING: //TC_STRING
		this.readTC_STRING()
		break

	case TC_CLASS: //TC_CLASS
		this.readNewClass()
		break

	case TC_ARRAY: //TC_ARRAY
		this.readNewArray()
		break

	case TC_ENUM: //TC_ENUM
		this.readNewEnum()
		break

	default: //Unknown/unsupported
		log.Panicln("Error: Unexpected identifier for object field value 0x" + this.byteToHex(this._data.peek()))
	}
	this.decreaseIndent()
}

/*******************
 * Read a double field.
 ******************/
func (this *SerializedObjectParser) readDoubleField() {
	var b1, b2, b3, b4, b5, b6, b7, b8 byte
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	b5 = this._data.pop()
	b6 = this._data.pop()
	b7 = this._data.pop()
	b8 = this._data.pop()
	var xx uint64 = (uint64(b1<<56) & 0xff00000000000000) +
		(uint64(b2<<48) & 0xff000000000000) +
		(uint64(b3<<40) & 0xff0000000000) +
		(uint64(b4<<32) & 0xff00000000) +
		(uint64(b5<<24) & 0xff000000) +
		(uint64(b6<<16) & 0xff0000) +
		(uint64(b7<<8) & 0xff00) +
		uint64(b8&0xff)
	this.print("(double)", xx, " - 0x"+this.byteToHex(b1)+
		" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4)+" "+this.byteToHex(b5)+" "+this.byteToHex(b6)+" "+
		this.byteToHex(b7)+" "+this.byteToHex(b8))
}

// 读新class
func (this *SerializedObjectParser) readNewClass() {
	var b1 byte

	//TC_CLASS
	b1 = this._data.pop()
	this.print("TC_CLASS - 0x" + this.byteToHex(b1))
	if b1 != TC_CLASS {
		log.Panicln("Error: Illegal value for TC_CLASS (should be 0x76)")
	}
	this.increaseIndent()

	//classDesc
	this.readClassDesc()

	//Revert indent
	this.decreaseIndent()

	//newHandle
	this.newHandle1()
}

// 读新数组
func (this *SerializedObjectParser) readNewArray() {
	var cdd *ClassDataDesc
	var cd *ClassDetails
	var b1, b2, b3, b4 byte
	var size int

	//TC_ARRAY
	b1 = this._data.pop()
	this.print("TC_ARRAY - 0x" + this.byteToHex(b1))
	if b1 != TC_ARRAY {
		log.Panicln("Error: Illegal value for TC_ARRAY (should be 0x75)")
	}
	this.increaseIndent()

	//classDesc
	cdd = this.readClassDesc() //Read the class data description to enable array elements to be read
	if cdd.getClassCount() != 1 {
		log.Panicln("Error: Array class description made up of more than one class.")
	}
	cd = cdd.getClassDetails(0)
	if cd.getClassName()[0:1] != "[" {
		log.Panicln("Error: Array class name does not begin with '['.")
	}

	//newHandle
	this.newHandle1()

	//Array size
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	size = int((uint32(b1<<24) & 0xff000000) +
		(uint32(b2<<16) & 0xff0000) +
		(uint32(b3<<8) & 0xff00) +
		(uint32(b4) & 0xff))
	this.print("Array size - ", size, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4))

	//Array data
	this.print("Values")
	this.increaseIndent()
	for i := 0; i < size; {
		i += 1
		//Print element index
		this.print("Index ", i, ":")
		this.increaseIndent()

		//Read the field values based on the classDesc read above
		x1 := []byte(cd.getClassName()[0:1])
		this.readFieldValue(x1[0])

		//Revert indent
		this.decreaseIndent()
	}
	this.decreaseIndent()

	//Revert indent
	this.decreaseIndent()
}
func (this *SerializedObjectParser) readTC_STRING() string {
	var val string
	var b1 byte

	// TC_STRING
	b1 = this._data.pop()
	this.print("TC_STRING - 0x" + this.byteToHex(b1))
	if b1 != TC_STRING {
		log.Panicln("Error: Illegal value for TC_STRING (should be 0x74)")
	}

	//Indent
	this.increaseIndent()

	//newHandle
	this.newHandle1()

	//UTF
	val = this.readUtf()

	//Decrease indent
	this.decreaseIndent()

	//Return the string value
	return val
}

func (this *SerializedObjectParser) readPrevObject() int {
	var b1, b2, b3, b4 byte
	var handle uint32

	//TC_REFERENCE
	b1 = this._data.pop()
	this.print("TC_REFERENCE - 0x" + this.byteToHex(b1))
	if b1 != TC_REFERENCE {
		log.Panicln("Error: Illegal value for TC_REFERENCE (should be 0x71)")
	}
	this.increaseIndent()

	//Reference handle
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()

	var a11 = []byte{b1, b2, b3, b4}
	handle = binary.BigEndian.Uint32(a11)
	//handle = (uint32(b1<<24) & 0xff000000) +
	//(uint32(b2<<16) & 0xff0000) +
	//(uint32(b3<<8) & 0xff00) +
	//(uint32(b4) & 0xff)

	this.print("Handle - ", handle, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4))

	//Revert indent
	this.decreaseIndent()

	//Return the handle
	return int(handle)
}

func (this *SerializedObjectParser) readNullReference() {
	var b1 byte

	//TC_NULL
	b1 = this._data.pop()
	this.print("TC_NULL - 0x" + this.byteToHex(b1))
	if b1 != TC_NULL {
		log.Panicln("Error: Illegal value for TC_NULL (should be 0x70)")
	}
}

func (this *SerializedObjectParser) readException() {}

func (this *SerializedObjectParser) handleReset() {}

func (this *SerializedObjectParser) readBlockData() {
	contents := ""
	var len int
	var b1 byte

	//TC_BLOCKDATA
	b1 = this._data.pop()
	this.print("TC_BLOCKDATA - 0x" + this.byteToHex(b1))
	if b1 != TC_BLOCKDATA {
		log.Panicln("Error: Illegal value for TC_BLOCKDATA (should be 0x77)")
	}
	this.increaseIndent()

	//size
	len = int(this._data.pop() & 0xFF)
	this.print("Length - ", len, " - 0x"+this.byteToHex((byte)(len&0xff)))

	//contents
	for i := 0; i < len; {
		i += 1
		contents += this.byteToHex(this._data.pop())
	}
	this.print("Contents - 0x" + contents)

	//Drop indent back
	this.decreaseIndent()
}

func (this *SerializedObjectParser) readLongBlockData() {
	contents := ""
	var len uint32
	var b1, b2, b3, b4 byte

	//TC_BLOCKDATALONG
	b1 = this._data.pop()
	this.print("TC_BLOCKDATALONG - 0x" + this.byteToHex(b1))
	if b1 != TC_BLOCKDATALONG {
		log.Panicln("Error: Illegal value for TC_BLOCKDATA (should be 0x77)")
	}
	this.increaseIndent()

	//size
	b1 = this._data.pop()
	b2 = this._data.pop()
	b3 = this._data.pop()
	b4 = this._data.pop()
	len = (uint32(b1<<24) & 0xff000000) +
		(uint32(b2<<16) & 0xff0000) +
		(uint32(b3<<8) & 0xff00) +
		(uint32(b4) & 0xff)
	this.print("Length - ", len, " - 0x"+this.byteToHex(b1)+" "+this.byteToHex(b2)+" "+this.byteToHex(b3)+" "+this.byteToHex(b4))

	//contents
	var l uint32 = 0
	for l < len {
		l += 1
		contents += this.byteToHex(this._data.pop())
	}
	this.print("Contents - 0x" + contents)

	//Drop indent back
	this.decreaseIndent()
}

// content reads the next object in the stream and parses it.
func (this *SerializedObjectParser) content(allowedNames map[string]bool) (content interface{}, err error) {
	var tc uint8

	tc = this._data.peek()
	this.so.Tc_Type = tc
	switch tc {
	case TC_NULL: // = 0x70 // 空指针
		this.readNullReference()
	case TC_REFERENCE: // = 0x71
		this.readPrevObject()
	case TC_CLASSDESC, TC_PROXYCLASSDESC: // = 0x7D TC_PROXYCLASSDESC: // = 0x72 // TC_CLASSDESC. 指定这是一个新类。
		this.readNewClassDesc()
	case TC_OBJECT: // = 0x73 // TC_OBJECT.  指定这是一个新的Object.
		this.readNewObject()
	case TC_STRING, TC_LONGSTRING: // = 0x7C: // = 0x74
		this.readNewString()
	case TC_ARRAY: // = 0x75
		this.readNewArray()
	case TC_CLASS: // = 0x76
		this.readNewClass()
	case TC_BLOCKDATA: // = 0x77
		this.readBlockData()
	case TC_ENDBLOCKDATA: // = 0x78
	case TC_RESET: // = 0x79
		this.handleReset()
	case TC_BLOCKDATALONG: // = 0x7A
		this.readLongBlockData()
	case TC_EXCEPTION: // = 0x7B
		this.readException()
	case TC_ENUM: // = 0x7E
		this.readNewEnum()
	default: // 异常情况
	}

	return nil, nil
}

// end check has next byte in stream.
func (this *SerializedObjectParser) end() bool {
	if this.rd.Buffered() == 0 {
		_, eof := this.rd.Peek(1)

		return eof != nil
	}

	return false
}

// readString reads a string of length cnt bytes.
func (this *SerializedObjectParser) readString(cnt int, asHex bool) (s string, err error) {
	this.buf.Reset()

	// Prevented to allocate an extremely large block of memory.
	if cnt > this.maxDataBlockSize {
		err = errors.Errorf("block data exceeds size of reader buffer. " +
			"To increase the size, use the method SetMaxDataBlockSize or use bufio.Reader with a larger buffer size")

		return
	}

	if _, err = io.CopyN(&this.buf, this.rd, int64(cnt)); err != nil {
		err = errors.Wrap(err, "error reading string")

		return
	}

	if asHex {
		s = hex.EncodeToString(this.buf.Bytes())
	} else {
		s = this.buf.String()
	}

	return
}

func (this *SerializedObjectParser) readUInt8() (x uint8, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading uint8")
	}

	return
}

func (this *SerializedObjectParser) readInt8() (x int8, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading int8")
	}

	return
}

func (this *SerializedObjectParser) readUInt16() (x uint16, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading uint16")
	}

	return
}

func (this *SerializedObjectParser) readInt16() (x int16, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading int16")
	}

	return
}

func (this *SerializedObjectParser) readUInt32() (x uint32, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading uint32")
	}

	return
}

func (this *SerializedObjectParser) readInt32() (x int32, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading int32")
	}

	return
}

func (this *SerializedObjectParser) readFloat32() (x float32, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading float32")
	}

	return
}

func (this *SerializedObjectParser) readInt64() (x int64, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading int64")
	}

	return
}

func (this *SerializedObjectParser) readFloat64() (x float64, err error) {
	if err = binary.Read(this.rd, binary.BigEndian, &x); err != nil {
		err = errors.Wrap(err, "error reading float64")
	}

	return
}

// utf reads a variable length string.
func (this *SerializedObjectParser) utf() (s string, err error) {
	var offset uint16

	if offset, err = this.readUInt16(); err != nil {
		err = errors.Wrap(err, "error reading utf: unable to read segment length")

		return
	}

	if s, err = this.readString(int(offset), false); err != nil {
		err = errors.Wrap(err, "error reading utf: unable to read segment")
	}

	return
}

// utf reads a large (up to 2^32 bytes) variable length string.
func (this *SerializedObjectParser) utfLong() (s string, err error) {
	var offset uint32

	if offset, err = this.readUInt32(); err != nil {
		err = errors.Wrap(err, "error reading utf: unable to read first segment length")

		return
	}

	if offset != 0 {
		err = errors.New("unable to read string larger than 2^32 bytes")

		return
	}

	if offset, err = this.readUInt32(); err != nil {
		err = errors.Wrap(err, "error reading utf long: unable to read second segment length")

		return
	}

	if s, err = this.readString(int(offset), false); err != nil {
		err = errors.Wrap(err, "error reading utf long: unable to read segment")
	}

	return
}

// magic checks for the presence of the STREAM_MAGIC value.
func (this *SerializedObjectParser) magic() error {
	magicVal, err := this.readUInt16()

	if err == nil && magicVal != STREAM_MAGIC {
		return errors.New("magic value STREAM_MAGIC not found")
	}
	this.so.STREAM_MAGIC = STREAM_MAGIC
	return err
}

// version checks to be sure the serialized object is using a supported protocol version.
func (this *SerializedObjectParser) version() error {
	ver, err := this.readUInt16()
	if err != nil {
		return err
	}

	if byte(ver) != STREAM_VERSION {
		return errors.Errorf("protocol version not recognized: wanted 5 got %d", ver)
	}
	this.so.STREAM_VERSION = STREAM_VERSION

	return nil
}

// field contains info about a single class member.
type field struct {
	className string
	typeName  string
	name      string
}

// fieldDesc reads a single field descriptor.
func (this *SerializedObjectParser) fieldDesc() (f *field, err error) {
	var typeDec uint8

	if typeDec, err = this.readUInt8(); err != nil {
		err = errors.Wrap(err, "error reading field type")

		return
	}

	var name string

	if name, err = this.utf(); err != nil {
		err = errors.Wrap(err, "error reading field name")

		return
	}

	typeName := string(typeDec)

	f = &field{
		typeName: typeName,
		name:     name,
	}

	if strings.Contains("[L", typeName) { //nolint
		var className interface{}

		if className, err = this.content(nil); err != nil {
			err = errors.Wrap(err, "error reading field class name")

			return
		}

		var isString bool
		if f.className, isString = className.(string); !isString {
			err = errors.New("unexpected field class name type")
		}
	}

	return
}

// annotations reads all class annotations.
func (this *SerializedObjectParser) annotations(allowedNames map[string]bool) (anns []interface{}, err error) {
	for {
		var ann interface{}

		if ann, err = this.content(allowedNames); err != nil {
			err = errors.Wrap(err, "error reading class annotation")

			return
		}

		if _, isEndBlock := ann.(endBlockT); isEndBlock {
			break
		}

		anns = append(anns, ann)
	}

	return
}

// clazz contains java class info.
type clazz struct {
	super            *clazz
	annotations      []interface{}
	fields           []*field
	serialVersionUID string
	name             string
	flags            uint8
	isEnum           bool
}

// classDesc reads a class descriptor.
func (this *SerializedObjectParser) classDesc() (cls *clazz, err error) {
	var x interface{}

	if x, err = this.content(allowedClazzNames); err != nil {
		err = errors.Wrap(err, "error reading class description")

		return
	}

	if x == nil {
		return
	}

	var isClazz bool
	if cls, isClazz = x.(*clazz); !isClazz {
		err = errors.New("unexpected type returned while reading class description")
	}

	return
}

// parseClassDesc parses a class descriptor.
//nolint:funlen
func parseClassDesc(this *SerializedObjectParser) (x interface{}, err error) {
	cls := &clazz{}

	if cls.name, err = this.utf(); err != nil {
		err = errors.Wrap(err, "error reading class name")

		return
	}

	const minClassNameLength = 2
	if len(cls.name) < minClassNameLength {
		err = errors.Wrapf(err, "invalid class name: '%s'", cls.name)

		return
	}

	const serialVersionUIDLength = 8
	if cls.serialVersionUID, err = this.readString(serialVersionUIDLength, true); err != nil {
		err = errors.Wrap(err, "error reading class serialVersionUID")

		return
	}

	this.newHandle(cls)

	if cls.flags, err = this.readUInt8(); err != nil {
		err = errors.Wrap(err, "error reading class flags")

		return
	}

	cls.isEnum = (cls.flags & 0x10) != 0

	var fieldCount uint16

	if fieldCount, err = this.readUInt16(); err != nil {
		err = errors.Wrap(err, "error reading class field count")

		return
	}

	for i := 0; i < int(fieldCount); i++ {
		var f *field

		if f, err = this.fieldDesc(); err != nil {
			err = errors.Wrap(err, "error reading class field")

			return
		}

		cls.fields = append(cls.fields, f)
	}

	if cls.annotations, err = this.annotations(nil); err != nil {
		err = errors.Wrap(err, "error reading class annotations")

		return
	}

	if cls.super, err = this.classDesc(); err != nil {
		err = errors.Wrap(err, "error reading class super")

		return
	}

	x = cls

	return
}

func parseClass(this *SerializedObjectParser) (cd interface{}, err error) {
	if cd, err = this.classDesc(); err != nil {
		err = errors.Wrap(err, "error parsing class")

		return
	}

	cd = this.newHandle(cd)

	return
}

func parseReference(this *SerializedObjectParser) (ref interface{}, err error) {
	var refIdx int32

	if refIdx, err = this.readInt32(); err != nil {
		err = errors.Wrap(err, "error reading reference index")

		return
	}

	const refIDMask = 0x7e0000
	i := int(refIdx - refIDMask)

	if i > -1 && i < len(this.handles) {
		ref = this.handles[i]
	}

	return
}

func parseArray(this *SerializedObjectParser) (arr interface{}, err error) {
	var cls *clazz

	if cls, err = this.classDesc(); err != nil {
		err = errors.Wrap(err, "error parsing array class")

		return
	}

	res := map[string]interface{}{
		"class": cls,
	}

	this.newHandle(res)

	var size int32

	if size, err = this.readInt32(); err != nil {
		err = errors.Wrap(err, "error reading array size")

		return
	}

	res["length"] = size

	if cls == nil {
		return
	}

	primHandler, exists := primitiveHandlers[string(cls.name[1])]
	if !exists {
		err = errors.Errorf("unknown field type '%s'", string(cls.name[1]))

		return
	}

	var array []interface{}

	for i := 0; i < int(size); i++ {
		var nxt interface{}

		if nxt, err = primHandler(this); err != nil {
			err = errors.Wrap(err, "error reading primitive array member")

			return
		}

		array = append(array, nxt)
	}

	arr = array

	return
}

// newDeferredHandle reserves an object handle slot and returns a func which can set the slot value at a later time.
func (this *SerializedObjectParser) newDeferredHandle() func(interface{}) interface{} {
	idx := len(this.handles)
	this.handles = append(this.handles, nil)

	return func(obj interface{}) interface{} {
		this.handles[idx] = obj

		return obj
	}
}

func parseEnum(this *SerializedObjectParser) (enum interface{}, err error) {
	var cls *clazz

	if cls, err = this.classDesc(); err != nil {
		err = errors.Wrap(err, "error parsing enum class")

		return
	}

	deferredHandle := this.newDeferredHandle()

	var enumConstant interface{}

	if enumConstant, err = this.content(nil); err != nil {
		err = errors.Wrap(err, "error parsing enum constant")

		return
	}

	res := map[string]interface{}{
		"value": enumConstant,
		"class": cls,
	}

	enum = deferredHandle(res)

	return
}

func parseBlockData(this *SerializedObjectParser) (bd interface{}, err error) {
	var size uint8

	if size, err = this.readUInt8(); err != nil {
		err = errors.Wrap(err, "error parsing block data size")

		return
	}

	data := make([]byte, size)

	if _, err = io.ReadFull(this.rd, data); err == nil {
		bd = data
	}

	return
}

func parseBlockDataLong(this *SerializedObjectParser) (bdl interface{}, err error) {
	var size uint32

	if size, err = this.readUInt32(); err != nil {
		err = errors.Wrap(err, "error parsing block data long size")

		return
	}

	// Prevented to allocate an extremely large block of memory.
	if int(size) > this.maxDataBlockSize {
		err = errors.Errorf("block data exceeds size of reader buffer. " +
			"To increase the size, use the method SetMaxDataBlockSize or use bufio.Reader with a larger buffer size")

		return
	}

	data := make([]byte, size)

	if _, err = io.ReadFull(this.rd, data); err == nil {
		bdl = data

		return
	}

	return
}

func parseString(this *SerializedObjectParser) (str interface{}, err error) {
	if str, err = this.utf(); err != nil {
		err = errors.Wrap(err, "error parsing string")
	} else {
		str = this.newHandle(str)
	}

	return
}

func parseLongString(this *SerializedObjectParser) (longStr interface{}, err error) {
	if longStr, err = this.utfLong(); err != nil {
		err = errors.Wrap(err, "error parsing long string")
	} else {
		this.newHandle(longStr)
	}

	return
}

func parseNull(_ *SerializedObjectParser) (interface{}, error) {
	return nil, nil
}

type endBlockT string

const endBlock endBlockT = "endBlock"

func parseEndBlockData(_ *SerializedObjectParser) (interface{}, error) {
	return endBlock, nil
}

// values reads primitive field values.
func (this *SerializedObjectParser) values(cls *clazz) (vals map[string]interface{}, err error) {
	var exists bool

	var handler primitiveHandler

	vals = make(map[string]interface{})

	for _, field := range cls.fields {
		if field == nil {
			continue
		}

		if handler, exists = primitiveHandlers[field.typeName]; !exists {
			err = errors.Errorf("unknown field type '%s'", field.typeName)

			return
		}

		if vals[field.name], err = handler(this); err != nil {
			err = errors.Wrap(err, "error reading primitive field value")

			return
		}
	}

	return
}

// annotationsAsMap reads values (when isBlock is false) and merges annotations then calls any relevant post processor.
func (this *SerializedObjectParser) annotationsAsMap(cls *clazz, isBlock bool) (data map[string]interface{}, err error) {
	if isBlock {
		data = make(map[string]interface{})
	} else if data, err = this.values(cls); err != nil {
		err = errors.Wrap(err, "error reading class data field values")

		return
	}

	var anns []interface{}

	if anns, err = this.annotations(nil); err != nil {
		err = errors.Wrap(err, "error reading annotations")

		return
	}

	data["@"] = anns

	if !isBlock {
		if postproc, exists := KnownPostProcs[cls.name+"@"+cls.serialVersionUID]; exists {
			data, err = postproc(data, anns)
		}
	}

	return
}

// classData reads a serialized class into a generic data structure.
func (this *SerializedObjectParser) classData(cls *clazz) (data map[string]interface{}, err error) {
	if cls == nil {
		return nil, errors.New("invalid class definition: nil")
	}

	const (
		ScSerializableWithoutWriteMethod = 0x02
		ScSerializableWithWriteMethod    = 0x03
		ScExternalizeWithBlockData       = 0x04
		ScExternalizeWithoutBlockData    = 0x0c
	)

	switch cls.flags & 0x0f {
	case ScSerializableWithoutWriteMethod: // SC_SERIALIZABLE without SC_WRITE_METHOD
		return this.values(cls)

	case ScSerializableWithWriteMethod: // SC_SERIALIZABLE with SC_WRITE_METHOD
		return this.annotationsAsMap(cls, false)

	case ScExternalizeWithBlockData: // SC_EXTERNALIZABLE without SC_BLOCKDATA
		return nil, errors.New("unable to parse version 1 external content")

	case ScExternalizeWithoutBlockData: // SC_EXTERNALIZABLE with SC_BLOCKDATA
		return this.annotationsAsMap(cls, true)

	default:
		return nil, errors.Errorf("unable to deserialize class with flags %#x", cls.flags)
	}
}

// recursiveClassData recursively reads inheritance tree until it reaches java.lang.object.
func (this *SerializedObjectParser) recursiveClassData(cls *clazz, obj map[string]interface{},
	seen map[*clazz]bool) error {
	if cls == nil {
		return nil
	}

	seen[cls] = true

	if cls.super != nil && !seen[cls.super] {
		seen[cls.super] = true
		if err := this.recursiveClassData(cls.super, obj, seen); err != nil {
			return err
		}
	}

	extends, isMap := obj["extends"].(map[string]interface{})
	if !isMap {
		return errors.New("unexpected extends value")
	}

	fields, err := this.classData(cls)
	if err != nil {
		return errors.Wrap(err, "error reading recursive class data")
	}

	extends[cls.name] = fields

	for name, val := range fields {
		obj[name] = val
	}

	return nil
}

func parseObject(this *SerializedObjectParser) (obj interface{}, err error) {
	var cls *clazz

	if cls, err = this.classDesc(); err != nil {
		err = errors.Wrap(err, "error reading object class")

		return
	}

	objMap := map[string]interface{}{
		"class":   cls,
		"extends": make(map[string]interface{}),
	}

	deferredHandle := this.newDeferredHandle()

	seen := map[*clazz]bool{}
	if err = this.recursiveClassData(cls, objMap, seen); err != nil {
		err = errors.Wrap(err, "error reading recursive class data")

		return
	}

	obj = deferredHandle(objMap)

	return
}

// postProcSize reads the object size as an int32 from the first data element.
func postProcSize(data []interface{}, offset int) (size int, err error) {
	if len(data) < 1 {
		err = errors.New("invalid data: at least one element required")

		return
	}

	b, isByteSlice := data[0].([]byte)
	if !isByteSlice {
		err = errors.New("unexpected data at position 0")

		return
	}

	const minLength = 4
	if len(b) < offset+minLength {
		err = errors.Errorf("incorrect data at position 0: wanted at least %d bytes, got %d", offset+minLength, len(b))

		return
	}

	var size32 int32
	if err = binary.Read(bytes.NewReader(b[offset:]), binary.BigEndian, &size32); err != nil {
		err = errors.Wrap(err, "error reading size")

		return
	}

	size = int(size32)

	return
}

// listPostProc populates the object value with a []interface{}.
func listPostProc(fields map[string]interface{}, data []interface{}) (map[string]interface{}, error) {
	size, err := postProcSize(data, 0)
	if err != nil {
		return nil, err
	}

	if len(data) != size+1 {
		return nil, errors.Errorf("incorrect number of elements: want %d got %d", size, len(data)-1)
	}

	if size > 1 {
		fields["value"] = data[1:size]
	} else {
		fields["value"] = make([]interface{}, 0)
	}

	return fields, err
}

// mapPostProc populates the object value with a map of key/value pairs.
func mapPostProc(fields map[string]interface{}, data []interface{}) (map[string]interface{}, error) {
	size, err := postProcSize(data, 4)
	if err != nil {
		return nil, err
	}

	if size*2+1 > len(data) {
		return nil, errors.Errorf("incorrect number of elements: want %d got %d", size, len(data)-1)
	}

	m := make(map[string]interface{})

	for i := 0; i < size; i++ {
		key := data[2*i+1]
		value := data[2*i+2]

		if s, isString := key.(string); isString {
			m[s] = value
		}
	}

	fields["value"] = m

	return fields, nil
}

// enumMapPostProc populates the object value with a map of key/value pairs where keys are enum constants.
func enumMapPostProc(fields map[string]interface{}, data []interface{}) (map[string]interface{}, error) {
	size, err := postProcSize(data, 0)
	if err != nil {
		return nil, err
	}

	if size*2+1 > len(data) {
		return nil, errors.Errorf("incorrect number of elements: want %d got %d", size, len(data)-1)
	}

	m := make(map[string]interface{})

	for i := 0; i < size; i++ {
		key := data[2*i+1]
		value := data[2*i+2]

		if mk, isMap := key.(map[string]interface{}); isMap {
			if s, isString := mk["value"].(string); isString {
				m[s] = value
			}
		}
	}

	fields["value"] = m

	return fields, nil
}

// hashSetPostProc populates the object value with a map of key/value pairs.
func hashSetPostProc(fields map[string]interface{}, data []interface{}) (map[string]interface{}, error) {
	size, err := postProcSize(data, 8)
	if err != nil {
		return nil, err
	}

	if len(data) != size+1 {
		return nil, errors.Errorf("incorrect number of elements: want %d got %d", size, len(data)-1)
	}

	m := make(map[string]bool)

	if size > 1 {
		for idx := range data[1:size] {
			key := data[idx+1]
			if s, isString := key.(string); isString {
				m[s] = true
			}
		}
	}

	fields["value"] = m

	return fields, nil
}

// datePostProc populates the object value with a time.Time.
func datePostProc(fields map[string]interface{}, data []interface{}) (map[string]interface{}, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid data: at least one element required")
	}

	b, isByteSlice := data[0].([]byte)
	if !isByteSlice {
		return nil, errors.New("unexpected data at position 0")
	}

	const timestampBlockSize = 8
	if len(b) < timestampBlockSize {
		return nil, errors.Errorf("incorrect data at position 0: wanted 8 bytes, got %d", len(b))
	}

	var timestamp int64
	if err := binary.Read(bytes.NewReader(b[0:timestampBlockSize]), binary.BigEndian, &timestamp); err != nil {
		return nil, errors.Wrap(err, "error reading timestamp")
	}

	fields["value"] = time.Unix(0, timestamp*int64(time.Millisecond))

	return fields, nil
}
