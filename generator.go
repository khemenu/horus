package horus

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

var (
	DefaultUserAliasGenerator   = NewStaticStringGenerator([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"), 9)
	DefaultOpaqueTokenGenerator = NewStaticStringGenerator([]rune("-.ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz0123456789~"), 64)
)

type Generator interface {
	New() (string, error)
}

type staticStringGenerator struct {
	letters []rune
	size    uint
}

func NewStaticStringGenerator(letters []rune, size uint) Generator {
	g := &staticStringGenerator{make([]rune, len(letters)), size}
	copy(g.letters, letters)
	return g
}

func (g *staticStringGenerator) New() (string, error) {
	rst := make([]rune, g.size)
	buff := make([]byte, 8)

	for i := range rst {
		if _, err := rand.Read(buff); err != nil {
			return "", fmt.Errorf("crypto rand: %w", err)
		}

		idx := binary.LittleEndian.Uint64(buff) % uint64(len(g.letters))
		rst[i] = g.letters[idx]
	}

	return string(rst), nil
}
