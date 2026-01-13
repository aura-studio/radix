package radix

import (
	"errors"
	"math"
	"math/big"
)

type Radix struct {
	*Charset
	endian  Endian
	base    int
	maxSize int // -1 means no limit
}

func New(endian Endian, maxSize int, c *Charset) *Radix {
	r := &Radix{
		Charset: c,
		endian:  endian,
		base:    c.Size(),
		maxSize: maxSize,
	}
	return r
}

func (r *Radix) EncodeBytes(data []byte) string {
	n := big.NewInt(0)
	for _, b := range data {
		n.Mul(n, big.NewInt(256))
		n.Add(n, big.NewInt(int64(b)))
	}
	return string(r.EncodeBigInt(n))
}

func (r *Radix) DecodeBytes(dataString string) ([]byte, error) {
	data := []byte(dataString)
	if r.endian == BigEndian {
		var copyData = make([]byte, len(data))
		copy(copyData, data)
		r.reverseEndian(copyData)
		data = copyData
	}

	var n big.Int
	for i := len(data) - 1; i >= 0; i-- {
		n.Mul(&n, big.NewInt(int64(r.base)))

		// check if the character is valid
		if _, ok := r.charsetMap[data[i]]; !ok {
			return nil, errors.New("invalid character")
		}

		n.Add(&n, big.NewInt(int64(r.charsetMap[data[i]])))
	}

	decodeData := make([]byte, 0)
	for n.Sign() > 0 {
		mod := new(big.Int).Mod(&n, big.NewInt(int64(256)))
		decodeData = append(decodeData, byte(mod.Int64()))
		n.Div(&n, big.NewInt(256))
	}

	r.reverseEndian(decodeData)
	return decodeData, nil
}

func (r *Radix) EncodeBigInt(n *big.Int) []byte {
	n = r.validateBigInt(n)
	var data []byte
	if r.maxSize > 0 {
		data = make([]byte, 0, r.maxSize)
	} else {
		data = make([]byte, 0)
	}
	for n.Sign() > 0 {
		mod := new(big.Int).Mod(n, big.NewInt(int64(r.base)))
		data = append(data, r.charset[mod.Int64()])
		n = new(big.Int).Div(n, big.NewInt(int64(r.base)))
	}
	if r.endian == BigEndian {
		r.reverseEndian(data)
	}
	return data
}

func (r *Radix) validateBigInt(n *big.Int) *big.Int {
	if r.maxSize > 0 {
		return n.Mod(n, big.NewInt(int64(math.Pow(float64(r.base), float64(r.maxSize)))))
	}
	return n
}

func (r *Radix) Encode(n uint64) []byte {
	n = r.validate(n)
	var data []byte
	if r.maxSize > 0 {
		data = make([]byte, 0, r.maxSize)
	} else {
		data = make([]byte, 0)
	}
	for n > 0 {
		data = append(data, r.charset[n%uint64(r.base)])
		n /= uint64(r.base)
	}
	if r.endian == BigEndian {
		r.reverseEndian(data)
	}
	return data
}

func (r *Radix) Decode(data []byte) uint64 {
	if r.endian == BigEndian {
		var copyData = make([]byte, len(data))
		copy(copyData, data)
		r.reverseEndian(copyData)
		data = copyData
	}

	var n uint64
	for i := len(data) - 1; i >= 0; i-- {
		n = n*uint64(r.base) + uint64(r.charsetMap[data[i]])
	}
	return r.validate(n)
}

func (r *Radix) validate(i uint64) uint64 {
	if r.maxSize > 0 {
		return i % uint64(math.Pow(float64(r.base), float64(r.maxSize)))
	}
	return i
}

func (r *Radix) reverseEndian(data []byte) {
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}
}
