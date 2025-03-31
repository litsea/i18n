package i18n

import (
	"golang.org/x/text/language"
)

type Option func(*I18n)

func WithLanguages(ls ...language.Tag) Option {
	return func(i *I18n) {
		if len(ls) > 0 {
			i.languages = ls
		}
	}
}

func WithDefaultLanguage(lng language.Tag) Option {
	return func(i *I18n) {
		i.defaultLanguage = lng
	}
}

func WithLoaders(ls ...Loader) Option {
	return func(i *I18n) {
		if len(ls) > 0 {
			i.loaders = ls
		}
	}
}

func WithLogger(l Logger) Option {
	return func(i *I18n) {
		i.logger = l
	}
}

func WithDefaultLogger() Option {
	return func(i *I18n) {
		i.logger = defaultLogger{}
	}
}

func (i *I18n) GetLanguages() []language.Tag {
	return i.languages
}

func (i *I18n) GetDefaultLanguage() language.Tag {
	return i.defaultLanguage
}
