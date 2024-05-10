package main

import (
	"fmt"
	"log"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

// CWD is assumed to be a project root.
func main() {
	proto_ext, err := entproto.NewExtension(
		entproto.SkipGenFile(),
		entproto.WithProtoDir("./proto"),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("new entproto extension: %w", err))
	}

	err = entc.Generate(
		"./schema",
		&gen.Config{
			Package: "khepri.dev/horus/ent",
			Target:  "./ent",
		},
		entc.Extensions(
			proto_ext,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}
