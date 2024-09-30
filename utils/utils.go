package utils

import (
	"bytes"
	"crypto/tls"
	"net"
	"reflect"
	"unsafe"
)

// ByteSliceToString is ...
func ByteSliceToString(b []byte) string {
	// return *(*string)(unsafe.Pointer(&b))
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToByteSlice is ...
func StringToByteSlice(s string) []byte {
	// return unsafe.Slice((*byte)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&s)))), len(s))
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func RewindConn(conn net.Conn, read []byte) net.Conn {
	if tlsConn, ok := conn.(*tls.Conn); ok {
		var (
			tlsInput, _ = reflect.TypeOf(tls.Conn{}).FieldByName("input")
			input       = (*bytes.Reader)(unsafe.Add(unsafe.Pointer(tlsConn), tlsInput.Offset))
			remaining   = input.Len()
			size        = int(input.Size())
			buffered    = len(read)
		)
		if buffered <= size {
			_, _ = input.Seek(0, 0) //reset the offset
		} else {
			buf := make([]byte, buffered+remaining)
			copy(buf, read)                   //copy the already read data into buf
			_, _ = input.Read(buf[buffered:]) //read remaining data append into buf
			input.Reset(buf)                  //reset the offset
		}
		return tlsConn
	} else {
		return NewRawConn(conn, read)
	}
}
