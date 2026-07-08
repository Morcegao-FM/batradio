// Package web embute o build do frontend React (frontend-web → npm run build
// gera os arquivos em web/dist).
package web

import "embed"

//go:embed all:dist
var Dist embed.FS
