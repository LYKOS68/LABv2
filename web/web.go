// Package web enthaelt alle Seiten, Vorlagen und statischen Dateien. Sie werden beim Bauen
// ins Programm eingebettet, damit der Server ohne zusaetzliche Dateien laeuft.
package web

import "embed"

//go:embed *.html static
var Files embed.FS
