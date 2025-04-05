package i18n

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/litsea/i18n/testdata"
)

type args struct {
	lng     language.Tag
	msgID   string
	tplData map[any]any
}

type testTranslation struct {
	name    string
	args    args
	want    string
	errFunc func(error) bool
}

func TestTranslate(t *testing.T) {
	t.Parallel()

	tests := []testTranslation{
		{
			name: "hello",
			args: args{
				lng:   language.English,
				msgID: "welcome",
			},
			want: "hello",
		},
		{
			name: "hello alex",
			args: args{
				lng:   language.English,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want: "hello alex",
		},
		{
			name: "18 years old",
			args: args{
				lng:   language.English,
				msgID: "welcomeWithAge",
				tplData: map[any]any{
					"age": "18",
				},
			},
			want: "I am 18 years old",
		},
		{
			name: "fallback to english",
			args: args{
				lng:   language.English,
				msgID: "fallbackToEnglish",
			},
			want: "fallback to english",
		},
		{
			name: "fallback to english2",
			args: args{
				lng:   language.English,
				msgID: "fallbackToEnglish2",
			},
			want: "fallback to english2",
		},
		{
			name: "fallback to msgID",
			args: args{
				lng:   language.English,
				msgID: "fallbackToMsgID",
			},
			want:    "fallbackToMsgID",
			errFunc: IsMessageNotFoundErr,
		},
		{
			name: "not exist msgID",
			args: args{
				lng:   language.English,
				msgID: "notExistMsgID",
			},
			want:    "notExistMsgID",
			errFunc: IsMessageNotFoundErr,
		},
		// German
		{
			name: "hallo",
			args: args{
				lng:   language.German,
				msgID: "welcome",
			},
			want: "hallo",
		},
		{
			name: "hallo alex",
			args: args{
				lng:   language.German,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want: "hallo alex",
		},
		{
			name: "18 jahre alt",
			args: args{
				lng:   language.German,
				msgID: "welcomeWithAge",
				tplData: map[any]any{
					"age": "18",
				},
			},
			want: "ich bin 18 Jahre alt",
		},
		{
			name: "fallback to english",
			args: args{
				lng:   language.German,
				msgID: "fallbackToEnglish",
			},
			want:    "fallback to english",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "fallback to english2",
			args: args{
				lng:   language.German,
				msgID: "fallbackToEnglish2",
			},
			want:    "fallback to english2",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "fallback to msgID",
			args: args{
				lng:   language.German,
				msgID: "fallbackToMsgID",
			},
			want:    "fallbackToMsgID",
			errFunc: IsMessageNotFoundErr,
		},
		{
			name: "not exist msgID",
			args: args{
				lng:   language.German,
				msgID: "notExistMsgID",
			},
			want:    "notExistMsgID",
			errFunc: IsMessageNotFoundErr,
		},
		// French (not exist language fallback)
		{
			name: "hello",
			args: args{
				lng:   language.French,
				msgID: "welcome",
			},
			want:    "hello",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "hello alex",
			args: args{
				lng:   language.French,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want:    "hello alex",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "18 years old",
			args: args{
				lng:   language.French,
				msgID: "welcomeWithAge",
				tplData: map[any]any{
					"age": "18",
				},
			},
			want:    "I am 18 years old",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "fallback to english",
			args: args{
				lng:   language.French,
				msgID: "fallbackToEnglish",
			},
			want:    "fallback to english",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "fallback to english2",
			args: args{
				lng:   language.French,
				msgID: "fallbackToEnglish2",
			},
			want:    "fallback to english2",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "fallback to msgID",
			args: args{
				lng:   language.French,
				msgID: "fallbackToMsgID",
			},
			want:    "fallbackToMsgID",
			errFunc: IsMessageNotFoundErr,
		},
		{
			name: "not exist msgID",
			args: args{
				lng:   language.French,
				msgID: "notExistMsgID",
			},
			want:    "notExistMsgID",
			errFunc: IsMessageNotFoundErr,
		},
	}

	for _, tt := range tests {
		name := "DefaultLoader-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name)

		name = "FileLoader-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(FileLoader(defaultLocalizeFileDir)),
		)

		name = "EmbedLoader-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(EmbedLoader(testdata.Localize, "./localize/")),
		)
	}
}

func TestTranslateOverride(t *testing.T) {
	t.Parallel()

	tests := []testTranslation{
		{
			name: "hello",
			args: args{
				lng:   language.English,
				msgID: "welcome",
			},
			want: "hello2",
		},
		{
			name: "origin hello alex",
			args: args{
				lng:   language.English,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want: "hello alex",
		},
		{
			name: "hello",
			args: args{
				lng:   language.German,
				msgID: "welcome",
			},
			want: "hallo2",
		},
		{
			name: "origin hello alex",
			args: args{
				lng:   language.German,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want: "hallo alex",
		},
		{
			name: "hello",
			args: args{
				lng:   language.French,
				msgID: "welcome",
			},
			want:    "hello2",
			errFunc: IsMessageFallbackToDefaultErr,
		},
		{
			name: "origin hello alex",
			args: args{
				lng:   language.French,
				msgID: "welcomeWithName",
				tplData: map[any]any{
					"name": "alex",
				},
			},
			want:    "hello alex",
			errFunc: IsMessageFallbackToDefaultErr,
		},
	}

	for _, tt := range tests {
		name := "FileLoader-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(
				FileLoader(defaultLocalizeFileDir),
				FileLoader("./testdata/localize-override/"),
			),
		)

		name = "EmbedLoader-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(
				EmbedLoader(testdata.Localize, "./localize/"),
				EmbedLoader(testdata.LocalizeOverride, "./localize-override/"),
			),
		)

		name = "File-Embed-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(
				FileLoader(defaultLocalizeFileDir),
				EmbedLoader(testdata.LocalizeOverride, "./localize-override/"),
			),
		)

		name = "Embed-File-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			WithLoaders(
				EmbedLoader(testdata.Localize, "./localize/"),
				FileLoader("./testdata/localize-override/"),
			),
		)
	}
}

func runTranslationTest(t *testing.T, tt testTranslation, name string, opts ...Option) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		ops := []Option{
			WithLanguages(tt.args.lng),
			WithDefaultLanguage(language.English),
		}

		if len(opts) > 0 {
			ops = append(ops, opts...)
		}

		i := New(ops...)

		got, err := i.Translate(tt.args.lng, tt.args.msgID, tt.args.tplData)
		if got != tt.want {
			assert.Equal(t, tt.want, got)
		}

		if tt.errFunc == nil {
			assert.NoError(t, err)
		} else if !tt.errFunc(err) {
			t.Errorf("TranslateWithConfig(),  msgID = %v, unexpected error: %v", tt.args.msgID, err)
		}
	})
}

func TestTranslateMessageLoadErrorLog(t *testing.T) {
	t.Parallel()

	type args struct {
		lng language.Tag
		ls  []language.Tag
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		notWant []string
	}{
		{
			name: "all-exists",
			args: args{
				lng: language.English,
				ls:  []language.Tag{language.German},
			},
			want: nil,
		},
		{
			name: "no-exists",
			args: args{
				lng: language.English,
				ls:  []language.Tag{language.French},
			},
			want: []string{
				"level=WARN",
				"lng=fr",
				"load i18n message file",
				ErrBundleLoadMessageFileFailed.Error(),
			},
			notWant: []string{
				"level=ERROR",
				"lng=en",
				"load i18n default message file",
			},
		},
		{
			name: "no-exists-default",
			args: args{
				lng: language.French,
				ls:  []language.Tag{language.English},
			},
			want: []string{
				"level=ERROR",
				"lng=fr",
				"load i18n default message file",
				ErrBundleLoadMessageFileFailed.Error(),
			},
			notWant: []string{
				"level=WARN",
				"lng=en",
				"load i18n message file",
			},
		},
		{
			name: "no-exists-all",
			args: args{
				lng: language.French,
				ls:  []language.Tag{language.English, language.Chinese},
			},
			want: []string{
				"level=WARN",
				"lng=zh",
				"load i18n message file",
				"level=ERROR",
				"lng=fr",
				"load i18n default message file",
				ErrBundleLoadMessageFileFailed.Error(),
			},
			notWant: []string{
				"lng=en",
			},
		},
	}

	for _, tt := range tests {
		name := "MessageLoadErrorLog-" + tt.name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := new(strings.Builder)
			l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))

			_ = New(
				WithLanguages(tt.args.ls...),
				WithDefaultLanguage(tt.args.lng),
				WithLoaders(EmbedLoader(testdata.Localize, "./localize/")),
				WithLogger(l),
			)

			if len(tt.want) == 0 {
				assert.Equal(t, "", buf.String())
			} else {
				for _, w := range tt.want {
					assert.Contains(t, buf.String(), w)
				}
			}

			if len(tt.notWant) > 0 {
				for _, w := range tt.notWant {
					assert.NotContains(t, buf.String(), w)
				}
			}

			buf.Reset()
		})
	}
}
