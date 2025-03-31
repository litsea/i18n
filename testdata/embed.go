package testdata

import (
	"embed"
)

//go:embed localize/*
var Localize embed.FS

//go:embed localize-override/*
var LocalizeOverride embed.FS
