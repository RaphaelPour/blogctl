package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {

	/* A non-special string should stay the same */
	require.Equal(t, "prolog", Slug("prolog"))

	/* Spaces should be replaced by hyphens */
	require.Equal(t, "debugging-nightmares", Slug("debugging nightmares"))

	/* Uppercase letters should get low */
	require.Equal(t, "news", Slug("News"))

	/* Special characters should get removed */
	require.Equal(t, "lisp", Slug("((lisp))"))

	/* Everything together */
	require.Equal(t, "2-be---2-be", Slug("2 BE || ! 2 BE"))
}
