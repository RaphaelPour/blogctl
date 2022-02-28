package highlighter

import (
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/stretchr/testify/require"
)

func TestRenderer(t *testing.T) {
	renderer := GetRenderer()
	require.IsType(t, &(html.Renderer{}), renderer)

	text := "```go\n" +
		"func main() {\n" +
		"	fmt.Println('Test %d', 3)\n" +
		"	// Comment\n" +
		"}\n" +
		"```"
	result := markdown.ToHTML([]byte(text), nil, renderer)
	expected := "<div class=\"highlight\"><pre><code><span class=\"keyword\">func</span> main() {\n\tfmt.Println(<span class=\"string\">&#39;Test %d&#39;</span>, <span class=\"number\">3</span>)\n\t<span class=\"comment\">// Comment</span>\n}\n</code></pre></div>"
	require.Equal(t, expected, string(result))
}
