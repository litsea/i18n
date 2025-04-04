package i18n

import (
	"embed"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	defaultLocalizeFileFormat = "yaml"
	defaultLocalizeFileDir    = "./testdata/localize/"
	defaultLoader             = FileLoader(defaultLocalizeFileDir)
)

type Loader interface {
	LoadMessage(bd *i18n.Bundle, lng language.Tag) error
}

type LoaderFunc func(lng language.Tag) (string, []byte, error)

func (f LoaderFunc) LoadMessage(bd *i18n.Bundle, lng language.Tag) error {
	src, buf, err := f(lng)
	if err != nil {
		return err
	}

	if _, err = bd.ParseMessageFileBytes(buf, src); err != nil {
		return err
	}

	return nil
}

func FileLoader(dir string) LoaderFunc {
	return func(lng language.Tag) (string, []byte, error) {
		src := getLanguageFilePath(dir, lng)
		buf, err := os.ReadFile(src)
		return src, buf, err
	}
}

func EmbedLoader(fs embed.FS, dir string) LoaderFunc {
	return func(lng language.Tag) (string, []byte, error) {
		src := getLanguageFilePath(dir, lng)
		buf, err := fs.ReadFile(src)
		return src, buf, err
	}
}

func (i *I18n) loadMessageFiles() {
	for _, lng := range i.languages {
		i.loadMessageFile(lng)
	}

	if !slices.Contains(i.languages, i.defaultLanguage) {
		i.loadMessageFile(i.defaultLanguage)
	}
}

func (i *I18n) loadMessageFile(lng language.Tag) {
	var (
		success  bool
		firstErr error
	)

	for _, l := range i.loaders {
		err := l.LoadMessage(i.bundle, lng)
		if err == nil {
			success = true
		} else if firstErr == nil {
			firstErr = err
		}
	}

	if !success && firstErr != nil && i.logger != nil {
		if lng != i.defaultLanguage {
			i.logger.Warn("load i18n message file", "lng", lng,
				"err", fmt.Errorf("%w, %w", ErrBundleLoadMessageFileFailed, firstErr))
		} else {
			i.logger.Error("load i18n default message file", "lng", lng,
				"err", fmt.Errorf("%w, %w", ErrBundleLoadMessageFileFailed, firstErr),
			)
		}
	}
}

func getLanguageFilePath(dir string, lng language.Tag) string {
	return path.Join(dir, lng.String()+"."+defaultLocalizeFileFormat)
}
