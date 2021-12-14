package languageConst

const (
	FrenchLanguage  = "fr"
	EnglishLanguage = "en"
)

func AllowedLanguages() []string {
	return []string{
		FrenchLanguage,
		EnglishLanguage,
	}
}
