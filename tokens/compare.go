package tokens

import (
	"bytes"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (k *Key) Compare(v []byte) error {
	var keyer Keyer
	switch s := k.State.(type) {
	case *Key_Argon2:
		keyer = s.Argon2
	default:
		return errors.New("unknown keyer")
	}

	k_, err := keyer.Key(v, WithSalt(k.Salt))
	if err != nil {
		return fmt.Errorf("keyer: %w", err)
	}

	if !bytes.Equal(k.Hash, k_.Hash) {
		return errors.New("hash mismatch")
	}

	return nil
}

func Compare(v []byte, h []byte) error {
	t := Key{}
	if err := proto.Unmarshal(h, &t); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return t.Compare(v)
}
