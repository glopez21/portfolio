package main

import "embed"

//go:embed templates/*.html static/*.css
var contentFS embed.FS

//go:embed all:templates/*.html all:static/*.css
var staticFS embed.FS
