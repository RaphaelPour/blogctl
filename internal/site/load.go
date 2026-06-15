package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/RaphaelPour/blogctl/internal/common"
	"github.com/RaphaelPour/blogctl/internal/highlighter"
	"github.com/RaphaelPour/blogctl/internal/metadata"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

const (
	indexFile        = "index.html"
	postFileTemplate = "%s.html"
)

// imageMacro matches the IMAGE(<filename>) shortcode used in post bodies.
var imageMacro = regexp.MustCompile(`IMAGE\(([\w\.]+)\)`)

// loadPosts discovers every post directory under BlogPath, loads and renders
// the public ones, sorts them newest-first and wires up navigation links.
func (s *Site) loadPosts() error {
	postDirs, err := os.ReadDir(s.Options.BlogPath)
	if err != nil {
		return fmt.Errorf("Error reading blog path: %s", err)
	}

	posts := make([]Post, 0)
	for i, dir := range postDirs {
		if !dir.IsDir() {
			continue
		}

		postPath := filepath.Join(s.Options.BlogPath, dir.Name())
		files, err := os.ReadDir(postPath)
		if err != nil {
			return fmt.Errorf("Error reading post path of %s: %s", postPath, err)
		}

		if len(files) < 2 {
			return fmt.Errorf(
				"Unexpected count of files in post path %s. Found: %d",
				postPath,
				len(files),
			)
		}

		meta, err := metadata.Load(postPath)
		if err != nil {
			return err
		}

		/* Overstep posts which aren't set to 'public' */
		if meta.Status != metadata.PUBLIC_STATUS {
			continue
		}

		fmt.Printf("Rendering post #%02d: %s\n", i, dir.Name())

		post, err := s.loadPost(postPath, meta)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}

	/* Sort posts newest-first */
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp > posts[j].Timestamp
	})

	linkNavigation(posts)

	s.Posts = posts
	return nil
}

// loadPost reads a single post's content, renders its markdown and resolves
// the IMAGE() shortcodes into <img> tags plus a list of images to copy.
func (s *Site) loadPost(postPath string, meta *metadata.Metadata) (Post, error) {
	slugTitle := common.Slug(meta.Title)

	content, err := os.ReadFile(common.GetContentFile(postPath))
	if err != nil {
		return Post{}, fmt.Errorf("Error reading post content %s: %s", postPath, err)
	}

	rendered := markdown.ToHTML(
		content, parser.NewWithExtensions(parser.CommonExtensions|parser.Footnotes),
		highlighter.GetRenderer(),
	)

	/* replace all IMAGE(<filename>) with valid path to filename */
	renderedStr := imageMacro.ReplaceAllString(
		string(rendered),
		fmt.Sprintf(`<img src="%s_$1"/>`, slugTitle),
	)

	images := make([]imageRef, 0)
	for _, match := range imageMacro.FindAllStringSubmatch(string(rendered), -1) {
		images = append(images, imageRef{
			src: filepath.Join(s.Options.BlogPath, fmt.Sprintf("%s/%s", slugTitle, match[1])),
			dst: filepath.Join(s.Options.OutPath, fmt.Sprintf("%s_%s", slugTitle, match[1])),
		})
	}

	var featuredImage string
	if len(meta.FeaturedImage) > 0 {
		featuredImage = fmt.Sprintf("https://%s/%s_%s", s.Config.Domain, slugTitle, meta.FeaturedImage)
	}

	return Post{
		Title:         meta.Title,
		Link:          fmt.Sprintf(postFileTemplate, slugTitle),
		PermaLink:     fmt.Sprintf("https://%s/%s.html", s.Config.Domain, slugTitle),
		HomeLink:      indexFile,
		Timestamp:     meta.CreatedAt,
		CreatedAt:     meta.Date(),
		Content:       renderedStr,
		Discussion:    s.Config.Discussion && !meta.Static,
		Rendered:      template.HTML(renderedStr),
		FeaturedImage: featuredImage,
		Metadata:      meta,
		images:        images,
	}, nil
}

// linkNavigation wires up previous/next links between non-static posts in the
// already-sorted slice.
func linkNavigation(posts []Post) {
	nextPost := -1
	for i := 0; i < len(posts); i++ {
		if posts[i].Metadata.Static {
			continue
		}

		if nextPost >= 0 {
			posts[i].NextPostLink = posts[nextPost].Link
			posts[nextPost].PreviousPostLink = posts[i].Link

			fmt.Println(posts[i].Title, "<->", posts[nextPost].Title)
		}
		nextPost = i
	}
}
