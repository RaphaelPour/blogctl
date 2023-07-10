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
	expected := "<pre tabindex=\"0\" class=\"chroma\"><code><span class=\"line\"><span class=\"cl\"><span class=\"kd\">func</span> <span class=\"nf\">main</span><span class=\"p\">(</span><span class=\"p\">)</span> <span class=\"p\">{</span>\n</span></span><span class=\"line\"><span class=\"cl\">\t<span class=\"nx\">fmt</span><span class=\"p\">.</span><span class=\"nf\">Println</span><span class=\"p\">(</span><span class=\"err\">&#39;</span><span class=\"nx\">Test</span> <span class=\"o\">%</span><span class=\"nx\">d</span><span class=\"err\">&#39;</span><span class=\"p\">,</span> <span class=\"mi\">3</span><span class=\"p\">)</span>\n</span></span><span class=\"line\"><span class=\"cl\">\t<span class=\"c1\">// Comment\n</span></span></span><span class=\"line\"><span class=\"cl\"><span class=\"c1\"></span><span class=\"p\">}</span>\n</span></span></code></pre>"
	require.Equal(t, expected, string(result))
}
