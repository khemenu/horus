package tokens

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"golang.org/x/crypto/argon2"
)

type Keyer interface {
	Key(v []byte) ([]byte, error)
	Compare(v []byte, h []byte) error
}

type Argon2iKeyerInit struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func NewArgon2iKeyer(init Argon2iKeyerInit) Keyer {
	return &argon2iKeyer{
		Argon2iKeyerInit: init,
	}
}

type argon2iState struct {
	time    uint32
	memory  uint32
	threads uint8
	salt    []byte
	key     []byte
}

type argon2iKeyer struct {
	Argon2iKeyerInit
}

func (k *argon2iKeyer) Key(v []byte) ([]byte, error) {
	s := make([]byte, 32)
	if _, err := rand.Read(s); err != nil {
		return nil, fmt.Errorf("rand: %w", err)
	}

	key := k.key(v, s)
	h, err := k.encode(key, s)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return h, nil
}

func (k *argon2iKeyer) key(v []byte, s []byte) []byte {
	return argon2.Key(v, s, k.Time, k.Memory, k.Threads, k.KeyLen)
}

func (k *argon2iKeyer) encode(key []byte, salt []byte) ([]byte, error) {
	if len(salt) > math.MaxUint16 {
		return nil, errors.New("salt too long")
	}

	var h bytes.Buffer
	vs := []any{
		uint8(argon2.Version),
		k.Time,
		k.Memory,
		k.Threads,
		uint16(len(salt)),
	}
	for _, v := range vs {
		if err := binary.Write(&h, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}

	if _, err := h.Write(salt); err != nil {
		return nil, err
	}
	if _, err := h.Write(key); err != nil {
		return nil, err
	}

	return h.Bytes(), nil
}

func (*argon2iKeyer) decode(h []byte) (argon2iState, error) {
	var rst argon2iState

	r := bytes.NewBuffer(h)
	version := uint8(0)
	salt_len := uint16(0)
	ps := []any{
		&version,     // 1
		&rst.time,    // 4
		&rst.memory,  // 4
		&rst.threads, // 1
		&salt_len,    // 2
	}
	for _, p := range ps {
		if err := binary.Read(r, binary.LittleEndian, p); err != nil {
			return rst, err
		}
	}

	if version != argon2.Version {
		return rst, errors.New("version incompatible")
	}

	const HeaderSize = 1 + 4 + 4 + 1 + 2
	if len(h) <= (HeaderSize + int(salt_len)) {
		return rst, errors.New("invalid length of hash")
	}

	rst.salt = h[HeaderSize : HeaderSize+salt_len]
	rst.key = h[HeaderSize+salt_len:]

	return rst, nil
}

func (k *argon2iKeyer) Compare(v []byte, hashed []byte) error {
	s, err := k.decode(hashed)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	// Do not fast error return by comparing key length
	// to avoid side-channel attack by time.

	// Keying `v` with same parameter of `hashed`.
	k_ := (&argon2iKeyer{Argon2iKeyerInit: Argon2iKeyerInit{
		Time:    s.time,
		Memory:  s.memory,
		Threads: s.threads,
		KeyLen:  uint32(len(s.key)),
	}}).key(v, s.salt)
	if !bytes.Equal(k_, s.key) {
		return errors.New("mismatch")
	}

	return nil
}
