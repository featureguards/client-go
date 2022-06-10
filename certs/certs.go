package certs

import (
	_ "embed"
)

//go:embed test.pem
var TestCA []byte
