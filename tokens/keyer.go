package tokens

type Keyer interface {
	Key(v []byte, opts ...KeyerOption) (*Token, error)
}

type KeyerOption func(opt *KeyerOptions)

type KeyerOptions struct {
	Salt []byte
}

func NewKeyerOptions(opts ...KeyerOption) KeyerOptions {
	o := KeyerOptions{}
	o.FromOpts(opts...)
	return o
}

func (o *KeyerOptions) FromOpts(opts ...KeyerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithSalt(v []byte) KeyerOption {
	return func(opt *KeyerOptions) {
		opt.Salt = v
	}
}
