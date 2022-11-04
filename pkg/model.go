package pkg

import (
	"bufio"
	"bytes"
)

// 流中子对象
type SerChildObj struct {
}

type SerObject struct {
	STREAM_MAGIC   uint16 `json:"STREAM_MAGIC"`
	STREAM_VERSION uint16 `json:"STREAM_VERSION"`
	Tc_Type        byte   `json:"tc_Type"`
}

type Smooth struct {
	_p *SerializedObjectParser
}

// SerializedObjectParser reads serialized java objects
// see: https://docs.oracle.com/javase/8/docs/platform/serialization/spec/protocol.html
type SerializedObjectParser struct {
	buf                    bytes.Buffer
	rd                     *bufio.Reader
	handles                []interface{}
	maxDataBlockSize       int
	_handleValue           int
	_indent                string
	_classDataDescriptions []*ClassDataDesc
	_data                  Smooth
	so                     *SerObject // 序列化对象
}

const bufferSize = 1024

type Option func(sop *SerializedObjectParser)