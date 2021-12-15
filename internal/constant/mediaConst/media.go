package mediaConst

const (
	DefaultLanguage  = FrenchLanguage
	DefaultMediaType = MovieType
)

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

const (
	MovieType = "movie"
	TvType    = "tv"
)

func AllowedMediaTypes() []string {
	return []string{
		MovieType,
		TvType,
	}
}
