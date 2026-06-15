package site

import (
	"html/template"

	"github.com/RaphaelPour/blogctl/internal/metadata"
)

// Post is the in-memory representation of a single blog post ready to be
// written to disk.
type Post struct {
	Title            string
	Link             string
	PermaLink        string
	PreviousPostLink string
	NextPostLink     string
	HomeLink         string
	Timestamp        int64
	CreatedAt        string
	Content          string
	FeaturedImage    string
	Discussion       bool
	Rendered         template.HTML
	Metadata         *metadata.Metadata

	// images lists files referenced via IMAGE() in the post body that need to
	// be copied from the post's source directory into the output directory.
	images []imageRef
}

// imageRef is a single image to copy from a post directory to the output dir.
type imageRef struct {
	src string
	dst string
}
