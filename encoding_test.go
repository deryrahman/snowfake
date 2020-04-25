package snowfake_test

import (
	"fmt"
	"testing"

	"github.com/deryrahman/snowfake"
)

func TestEncodeBase58(t *testing.T) {
	name := func(id uint64) string {
		return fmt.Sprintf("id=%d", id)
	}
	tests := []struct {
		id              uint64
		expectedEncoded string
	}{
		{
			id:              uint64(18),
			expectedEncoded: "j",
		},
		{
			id:              uint64(1577836800),
			expectedEncoded: "3pqPUL",
		},
		{
			id:              0xFFFFFFFFFFFFFFFF,
			expectedEncoded: "JPwcyDCgEup",
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.id), func(t *testing.T) {
			encoded := snowfake.EncodeBase58(tt.id)
			assertEqual(t, tt.expectedEncoded, encoded)
		})
	}
}

func TestDecodeBase58(t *testing.T) {
	name := func(str string) string {
		return fmt.Sprintf("str=%s", str)
	}
	tests := []struct {
		str             string
		expectedDecoded uint64
	}{
		{
			str:             "j",
			expectedDecoded: uint64(18),
		},
		{
			str:             "3pqPUL",
			expectedDecoded: uint64(1577836800),
		},
		{
			str:             "JPwcyDCgEup",
			expectedDecoded: 0xFFFFFFFFFFFFFFFF,
		},
		{
			str:             "3pqPU=",
			expectedDecoded: 0,
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.str), func(t *testing.T) {
			decoded := snowfake.DecodeBase58(tt.str)
			assertEqual(t, tt.expectedDecoded, decoded)
		})
	}
}
