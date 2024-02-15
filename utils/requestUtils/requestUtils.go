package requestUtils

import (
	"bytes"
	"fmt"
	"github.com/ValeryVerkhoturov/chat/utils/i18nUtils"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"os"
	"strings"
)

func WrapHTMLWithEmbeddingJS(buf bytes.Buffer) string {
	return fmt.Sprintf(`
(function() {
	var wrapper = document.createElement("div");
	wrapper.innerHTML = unescape(`+"`%s`"+`);
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
})();
        `, &buf)
}

func GetLocale(c *gin.Context) (i18nUtils.Locale, string) {
	localeName := "ru"
	lang := c.Query("lang")

	locale, ok := i18nUtils.LocalesMap[lang]
	if ok {
		localeName = lang
	} else {
		locale = i18nUtils.LocalesMap[localeName]
	}
	return locale, localeName
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
