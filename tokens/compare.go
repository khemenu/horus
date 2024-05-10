package tokens

import (
	"bytes"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (t *Token) Compare(v []byte) error {
	var keyer Keyer
	switch s := t.State.(type) {
	case *Token_Argon2:
		keyer = s.Argon2
	default:
		return errors.New("unknown keyer")
	}

	t_, err := keyer.Key(v, WithSalt(t.Salt))
	if err != nil {
		return fmt.Errorf("keyer: %w", err)
	}

	if !bytes.Equal(t.Hash, t_.Hash) {
		return errors.New("hash mismatch")
	}

	return nil
}

func Compare(v []byte, h []byte) error {
	t := Token{}
	if err := proto.Unmarshal(h, &t); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return t.Compare(v)
}
