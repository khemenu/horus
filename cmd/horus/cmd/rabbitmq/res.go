package rabbitmq

import "io"

func Deny(w io.Writer) {
	w.Write([]byte("deny"))
}

func Allow(w io.Writer) {
	w.Write([]byte("allow"))
}
