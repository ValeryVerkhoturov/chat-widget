package requestUtils

import (
	"bytes"
	"fmt"
	"github.com/ValeryVerkhoturov/chat/utils/i18nUtils"
	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"html/template"
	"io/fs"
	"os"
	"strings"
)

type mediaType string

const (
	HTML mediaType = "text/html"
	SVG  mediaType = "image/svg+xml"
	JS   mediaType = "application/javascript"
)

const wrapTemplate = `(function() {
	var wrapper = document.createElement("div");
	wrapper.innerHTML = unescape(` + "`%s`" + `);
	document.body.appendChild(wrapper);

	// Move all scripts from wrapper to real script elements
	Array.from(wrapper.querySelectorAll('script')).forEach(function(oldScript) {
		var newScript = document.createElement('script');

		Array.from(oldScript.attributes).forEach(function(attr) {
			newScript.setAttribute(attr.name, attr.value);
		});

		if (oldScript.src) {
			newScript.src = oldScript.src;
		} else {
			newScript.textContent = oldScript.textContent;
		}
	
		oldScript.parentNode.replaceChild(newScript, oldScript);
	});
})();`

func createMinifier() *minify.M {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("application/javascript", js.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	return minifier
}

func minifyMediaType[T bytes.Buffer | []byte | string](minifier *minify.M, value T, mediaType mediaType) ([]byte, error) {
	switch v := any(value).(type) {
	case bytes.Buffer:
		return minifier.Bytes(string(mediaType), v.Bytes())
	case []byte:
		return minifier.Bytes(string(mediaType), v)
	case string:
		s, err := minifier.String(string(mediaType), v)
		return []byte(s), err
	default:
		return nil, fmt.Errorf("unknown type %T to minify as %s", v, mediaType)
	}
}

func GetWrappedHTMLWithEmbeddingJS(buf bytes.Buffer) (string, error) {
	minifier := createMinifier()

	minifiedHTML, err := minifyMediaType(minifier, buf, HTML)
	if err != nil {
		return "", err
	}

	jsScript := fmt.Sprintf(wrapTemplate, minifiedHTML)

	minifiedJsScript, err := minifyMediaType(minifier, jsScript, JS)
	if err != nil {
		return "", err
	}

	return string(minifiedJsScript), nil
}

func GetLocaleWithCode(c *gin.Context) (i18nUtils.Locale, string) {
	localeCode := "ru" // default
	lang := c.Query("lang")

	locale, ok := i18nUtils.LocalesMap[lang]
	if ok {
		localeCode = lang
	} else {
		locale = i18nUtils.LocalesMap[localeCode]
	}
	return locale, localeCode
}

// TemplateParseFSRecursive recursively parses all templates in the FS with the given extension.
// File paths are used as template names to support duplicate file names.
// Use nonRootTemplateNames to exclude root directory from template names
// (e.g. index.html instead of templates/index.html)
func TemplateParseFSRecursive(
	templates fs.FS,
	ext string,
	nonRootTemplateNames bool,
	funcMap template.FuncMap) (*template.Template, error) {

	root := template.New("")
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}
			b, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}
			name := ""
			if nonRootTemplateNames {
				//name the template based on the file path (excluding the root)
				parts := strings.Split(path, string(os.PathSeparator))
				name = strings.Join(parts[1:], string(os.PathSeparator))
			}
			t := root.New(name).Funcs(funcMap)
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return root, err
}

func GetApplicationError(errorDescription string) gin.H {
	return gin.H{"error": errorDescription}
}
