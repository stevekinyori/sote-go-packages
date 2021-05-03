package sHelper

import (
	"fmt"
	"testing"
	"time"
)

func TestLibUUID(t *testing.T) {
	AssertEqual(t, UUID(UUIDKind.Time)[:3], fmt.Sprintf("%v", time.Millisecond)[:3])
	AssertEqual(t, len(UUID(UUIDKind.Short)), 13)
	AssertEqual(t, len(UUID(UUIDKind.Long)), 36)
}
