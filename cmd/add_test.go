package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {

	/* A non-special string should stay the same */
	require.Equal(t, "prolog", slug("prolog"))

	/* Spaces should be replaced by hyphens */
	require.Equal(t, "debugging-nightmares", slug("debugging nightmares"))

	/* Uppercase letters should get low */
	require.Equal(t, "news", slug("News"))

	/* Special characters should get removed */
	require.Equal(t, "lisp", slug("((lisp))"))

	/* Everything together */
	require.Equal(t, "2-be--2-be", slug("2 BE || ! 2 BE"))
}
