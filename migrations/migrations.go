package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var embedFS embed.FS

var FS, _ = fs.Sub(embedFS, ".")
