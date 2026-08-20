package site

import (
	"embed"
	"html/template"
	"path"
	plainTemplate "text/template"

	"github.com/RaphaelPour/blogctl/internal/common"
)

const publicDir = "public"

var (
	//go:embed public/*
	publicFS embed.FS

	// postTmpl and staticTmpl render individual post pages with html/template.
	postTmpl   = template.Must(template.New("post").Parse(mustReadAsset("post.tmpl.html")))
	staticTmpl = template.Must(template.New("static").Parse(mustReadAsset("static.tmpl.html")))

	// indexTmpl renders the start page. It deliberately uses text/template to
	// preserve the original (unescaped) rendering behaviour.
	indexTmpl = plainTemplate.Must(plainTemplate.New("blog").Parse(mustReadAsset("index.tmpl.html")))
)

// mustReadAsset reads a file embedded from the default theme. The assets are
// guaranteed to exist at build time, so a failure here is a programmer error.
func mustReadAsset(name string) string {
	return string(common.Unwrap(publicFS.ReadFile(path.Join(publicDir, name))))
}
