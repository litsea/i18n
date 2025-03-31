# i18n

## Usage

### Initialize

#### File Loader

```golang
import (
	"gitlab.com/frontier3302117/golang/i18n"
	"golang.org/x/text/language"
)

i18 := i18n.New(
	i18n.WithLanguages([]language.Tag{language.English, language.German}),
	i18n.WithDefaultLanguage(language.English),
	i18n.WithLoaders([]i18n.Loader{
		FileLoader("./testdata/localize/"),
	}),
)
```

#### Embed Loader

```golang
import (
	"gitlab.com/frontier3302117/golang/i18n"
	"gitlab.com/frontier3302117/golang/i18n/testdata"
	"golang.org/x/text/language"
)

i18 := i18n.New(
	i18n.WithLanguages([]language.Tag{language.English, language.German, language.French}),
	i18n.WithDefaultLanguage(language.English),
	i18n.WithLoaders([]i18n.Loader{
		EmbedLoader(testdata.Localize, "./localize/"),
	}),
)
```

### Translate

```golang
import (
	"golang.org/x/text/language"
)

var msg string

// hello
msg = i18.Translate(language.English, "welcome")

// hallo alex
msg = i18.Translate(language.German, "welcomeWithName", map[any]any{
	"name": "alex",
})

// hello alex (fallback to English)
msg = i18.Translate(language.French, "welcomeWithName", map[any]any{
	"name": "alex",
})
```