package sHelper

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

type uuidKind int

type kinds struct {
	Time  uuidKind
	Short uuidKind
	Long  uuidKind
}

var UUIDKind = &kinds{
	Time:  0,
	Short: 1,
	Long:  2,
}

func UUID(kind uuidKind) string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if kind == UUIDKind.Time || n != len(uuid) || err != nil {
		return fmt.Sprintf("%v", time.Millisecond)
	}
	if kind == UUIDKind.Short {
		return fmt.Sprintf("%x-%x", uuid[0:4], uuid[4:6])
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
