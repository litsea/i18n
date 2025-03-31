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
		if err := i.loadMessageFile(lng); err != nil && i.logger != nil {
			if lng != i.defaultLanguage {
				i.logger.Warn(fmt.Errorf("load message file: %w, %w",
					ErrBundleLoadMessageFileFailed, err).Error(),
					"lng", lng,
				)
			} else {
				i.logger.Error(fmt.Errorf("load default message: %w, %w",
					ErrBundleLoadMessageFileFailed, err).Error(),
					"lng", lng,
				)
			}
		}
	}

	if !slices.Contains(i.languages, i.defaultLanguage) {
		if err := i.loadMessageFile(i.defaultLanguage); err != nil && i.logger != nil {
			i.logger.Error(fmt.Errorf("load default message: %w, %w",
				ErrBundleLoadMessageFileFailed, err).Error(),
				"lng", i.defaultLanguage,
			)
		}
	}
}

func (i *I18n) loadMessageFile(lng language.Tag) error {
	for _, l := range i.loaders {
		err := l.LoadMessage(i.bundle, lng)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLanguageFilePath(dir string, lng language.Tag) string {
	return path.Join(dir, lng.String()+"."+defaultLocalizeFileFormat)
}
