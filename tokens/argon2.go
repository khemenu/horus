package tokens

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"google.golang.org/protobuf/proto"
)

func NewArgon2i(s *Argon2State) Keyer {
	s.Version = argon2.Version
	s.HashType = 1
	return s
}

func (s *Argon2State) Key(v []byte, opts ...KeyerOption) (*Key, error) {
	if s.Version != argon2.Version {
		return nil, errors.New("version mismatch")
	}

	t := &Key{State: &Key_Argon2{Argon2: proto.Clone(s).(*Argon2State)}}
	o := NewKeyerOptions(opts...)

	t.Salt = o.Salt
	if t.Salt == nil {
		t.Salt = make([]byte, 32)
		if _, err := rand.Read(t.Salt); err != nil {
			return nil, fmt.Errorf("rand: %w", err)
		}
	}

	t.Hash = s.key(v, t.Salt)
	return t, nil
}

func (s *Argon2State) key(v []byte, salt []byte) []byte {
	return argon2.Key(v, salt, s.Iterations, s.MemorySize, uint8(s.Parallelism), s.TagLength)
}
