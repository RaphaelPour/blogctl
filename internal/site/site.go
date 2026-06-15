package site

import (
	"fmt"
	"os"
	"time"

	"github.com/RaphaelPour/blogctl/internal/config"

	"github.com/gorilla/feeds"
)

// Options configures a single render run.
type Options struct {
	BlogPath string
	OutPath  string
	Force    bool
}

// Site is the in-memory model of a blog ready to be rendered. It owns the
// loaded configuration, the loaded posts and the feed accumulated while
// rendering. Build stages share a *Site, so a later stage can rely on state
// produced by an earlier one.
type Site struct {
	Config    *config.Config
	Options   Options
	Posts     []Post // all public posts, sorted newest-first, including static pages
	Published []Post // non-static posts shown on the index and in the feed
	Feed      *feeds.Feed
}

// New loads the blog configuration and all public posts from disk, rendering
// their markdown and wiring up next/previous navigation. It performs no writes.
func New(opts Options) (*Site, error) {
	cfg, err := config.Load(opts.BlogPath)
	if err != nil {
		return nil, err
	}

	s := &Site{
		Config:  cfg,
		Options: opts,
		Feed: &feeds.Feed{
			Title:       cfg.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s", cfg.Domain)},
			Description: cfg.Description,
			Author:      &feeds.Author{Name: cfg.Author},
			Created:     time.Now(),
			Items:       make([]*feeds.Item, 0),
		},
	}

	if err := s.loadPosts(); err != nil {
		return nil, err
	}

	return s, nil
}

// Render writes the static website to the configured output directory by
// running each build stage in order. New output artifacts (sitemaps, tag
// pages, ...) can be added by implementing Stage and appending it in stages().
func (s *Site) Render() error {
	if err := s.prepareOutput(); err != nil {
		return err
	}

	for _, stage := range s.stages() {
		if err := stage.Run(s); err != nil {
			return fmt.Errorf("%s: %w", stage.Name(), err)
		}
	}

	return nil
}

func (s *Site) prepareOutput() error {
	if _, err := os.Stat(s.Options.OutPath); !os.IsNotExist(err) && !s.Options.Force {
		return fmt.Errorf("Output folder already exists")
	}

	if err := os.MkdirAll(s.Options.OutPath, os.ModePerm); err != nil {
		return fmt.Errorf("Error creating output folder: %s", err)
	}

	return nil
}
