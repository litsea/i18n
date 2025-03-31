package i18n

import (
	"testing"

	"github.com/litsea/i18n/testdata"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
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
			FileLoader(defaultLocalizeFileDir),
		)

		name = "EmbedLoader-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			EmbedLoader(testdata.Localize, "./localize/"),
		)
	}

	overrideTests := []testTranslation{
		{
			name: "hello",
			args: args{
				lng:   language.English,
				msgID: "welcome",
			},
			want: "hello2",
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
			name: "hello",
			args: args{
				lng:   language.French,
				msgID: "welcome",
			},
			want:    "hello2",
			errFunc: IsMessageFallbackToDefaultErr,
		},
	}

	for _, tt := range overrideTests {
		name := "FileLoader-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			FileLoader(defaultLocalizeFileDir),
			FileLoader("./testdata/localize-override/"),
		)

		name = "EmbedLoader-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			EmbedLoader(testdata.Localize, "./localize/"),
			EmbedLoader(testdata.LocalizeOverride, "./localize-override/"),
		)

		name = "File-Embed-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			FileLoader(defaultLocalizeFileDir),
			EmbedLoader(testdata.LocalizeOverride, "./localize-override/"),
		)

		name = "Embed-File-Override-" + tt.args.lng.String() + "-" + tt.args.msgID + "-" + tt.name
		runTranslationTest(t, tt, name,
			EmbedLoader(testdata.Localize, "./localize/"),
			FileLoader("./testdata/localize-override/"),
		)
	}
}

func runTranslationTest(t *testing.T, tt testTranslation, name string, ls ...Loader) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		i := New(
			WithLanguages(tt.args.lng),
			WithDefaultLanguage(language.English),
			WithLoaders(ls...),
		)

		got, err := i.Translate(tt.args.lng, tt.args.msgID, tt.args.tplData)
		if got != tt.want {
			assert.Equal(t, tt.want, got)
		}

		if tt.errFunc == nil {
			assert.NoError(t, err)
		} else {
			if !tt.errFunc(err) {
				t.Errorf("TranslateWithConfig(),  msgID = %v, unexpected error: %v", tt.args.msgID, err)
			}
		}
	})
}
