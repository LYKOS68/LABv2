// Package web baut die Seiten und Dateien aus diesem Ordner ins Programm ein:
// studio.html und connect.html als Seiten, static/ als Dateien unter /static/.
package web

import "embed"

//go:embed studio.html connect.html static
var Files embed.FS
