package assets

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type EnumsLoader struct {
	asset    AssetFn
	assetDir AssetDirFn
	data     map[string]interface{}
}

func NewEnumsLoader() *EnumsLoader {
	return &EnumsLoader{
		asset:    Asset,
		assetDir: AssetDir,
		data:     map[string]interface{}{},
	}
}

func (l *EnumsLoader) Data() map[string]interface{} {
	return l.data
}

func (l *EnumsLoader) loadDir(dir string) error {
	files, err := l.assetDir(dir)
	if err != nil {
		return err
	}

	for _, fp := range files {
		looksLikeEnum := strings.HasSuffix(fp, ".json")
		if !looksLikeEnum {
			l.loadDir(path.Join(dir, fp))
			continue
		}

		name := path.Join(dir, fp)
		bytes, err := l.asset(name)
		if err != nil {
			return err
		}
		key := strings.TrimSuffix(fp, ".json")
		var value interface{}
		if err = json.Unmarshal(bytes, &value); err != nil {
			return errors.Wrapf(err, "failed to unmarshal '%s'", key)
		}
		l.data[key] = value
	}
	return nil
}
