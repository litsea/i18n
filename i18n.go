package i18n

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	ErrBundleLoadMessageFileFailed = errors.New("failed to load bundle message file")
	ErrMessageFallbackToDefault    = errors.New("message fallback to default language")
)

type I18n struct {
	bundle          *i18n.Bundle
	languages       []language.Tag
	defaultLanguage language.Tag
	loaders         []Loader
	localizerByLng  map[language.Tag]*i18n.Localizer
	logger          Logger
}

func New(opts ...Option) *I18n {
	i := &I18n{
		languages:       []language.Tag{language.English},
		defaultLanguage: language.English,
		loaders:         []Loader{defaultLoader},
		localizerByLng:  make(map[language.Tag]*i18n.Localizer),
	}

	for _, opt := range opts {
		opt(i)
	}

	bundle := i18n.NewBundle(i.defaultLanguage)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	i.bundle = bundle
	i.loadMessageFiles()

	for _, l := range i.languages {
		i.localizerByLng[l] = i.newLocalizer(l)
	}

	// Add default language if it isn't exist
	if _, hasDefaultLng := i.localizerByLng[i.defaultLanguage]; !hasDefaultLng {
		i.localizerByLng[i.defaultLanguage] = i.newLocalizer(i.defaultLanguage)
	}

	return i
}

func (i *I18n) newLocalizer(lng language.Tag) *i18n.Localizer {
	lngDefault := i.defaultLanguage.String()
	ls := []string{
		lng.String(),
	}

	if lng.String() != lngDefault {
		ls = append(ls, lngDefault)
	}

	localizer := i18n.NewLocalizer(
		i.bundle,
		ls...,
	)

	return localizer
}

// Translate message.
func (i *I18n) Translate(lng language.Tag, msgID string, tplData ...map[any]any) (string, error) {
	localizer, ok := i.localizerByLng[lng]
	if !ok {
		localizer = i.localizerByLng[i.defaultLanguage]
	}

	cfg := &i18n.LocalizeConfig{
		MessageID: msgID,
	}

	if len(tplData) > 0 && tplData[0] != nil {
		cfg.TemplateData = tplData[0]
	}

	msg, l, err := localizer.LocalizeWithTag(cfg)
	if err != nil {
		// Fallback to default language
		if !l.IsRoot() {
			return msg, ErrMessageFallbackToDefault
		}

		// Fallback to language.Und
		return msgID, err
	}

	if l != lng {
		return msg, ErrMessageFallbackToDefault
	}

	return msg, nil
}

func IsMessageNotFoundErr(err error) bool {
	var notFoundErr *i18n.MessageNotFoundErr
	ok := errors.As(err, &notFoundErr)

	return ok
}

func IsMessageFallbackToDefaultErr(err error) bool {
	return errors.Is(err, ErrMessageFallbackToDefault)
}
