package radix_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/aura-studio/radix"
)

func TestRadix(t *testing.T) {
	c := radix.NewCharset([]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")).Shuffle(0)
	if c.Size() != 64 {
		t.Error("radix charset size error")
	}

	r := radix.New(radix.BigEndian, 8, c)
	t.Log(r.Encode(123456789))

	if r.Decode(r.Encode(123456789)) != 123456789 {
		t.Error("radix decode error")
	}

	r2 := radix.New(radix.LittleEndian, 8, c)
	t.Log(r2.Encode(123456789))

	if r2.Decode(r2.Encode(123456789)) != 123456789 {
		t.Error("radix decode error")
	}
}

// generate a random string of 62 characters
func generateRandomString(length int) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestRadixBytes(t *testing.T) {
	c := radix.NewCharset([]byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")).Shuffle(20250411)
	if c.Size() != 62 {
		t.Error("radix charset size error")
	}

	r := radix.New(radix.BigEndian, -1, c)

	data := []byte(generateRandomString(256))
	t.Log(r.EncodeBytes(data))

	decodedBytes, err := r.DecodeBytes(r.EncodeBytes(data))
	if err != nil {
		t.Error("radix decode error")
	}

	if !slices.Equal(decodedBytes, data) {
		t.Error("radix decode error")
	}

	r2 := radix.New(radix.LittleEndian, -1, c)
	t.Log(r2.EncodeBytes(data))

	decodedBytes, err = r2.DecodeBytes(r2.EncodeBytes(data))
	if err != nil {
		t.Error("radix decode error")
	}

	if !slices.Equal(decodedBytes, data) {
		t.Error("radix decode error")
	}
}
