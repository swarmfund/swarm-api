package assets

import (
	"fmt"
	"html/template"
	"path"
	"strings"
)

type TemplatesLoader struct {
	asset    AssetFn
	assetDir AssetDirFn
	template *template.Template
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{
		asset:    Asset,
		assetDir: AssetDir,
		template: template.New("template"),
	}
}

func (t *TemplatesLoader) loadDir(dir string) error {
	files, err := t.assetDir(dir)
	if err != nil {
		return err
	}

	for _, fp := range files {
		looksLikeTemplate := strings.HasSuffix(fp, ".html")
		if !looksLikeTemplate {
			t.loadDir(path.Join(dir, fp))
			continue
		}
		name := path.Join(dir, fp)
		bytes, err := t.asset(name)
		if err != nil {
			return err
		}
		_, err = t.template.New(name).Parse(string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplatesLoader) Lookup(name string) *template.Template {
	name = fmt.Sprintf("templates/%s.html", name)
	return t.template.Lookup(name)
}
