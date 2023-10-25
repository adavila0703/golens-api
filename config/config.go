package config

type config struct {
	AllowApiKey string `split_words:"true" required:"true"`
	AllowOrigin string `split_words:"true" required:"true"`
	DatabaseURL string `split_words:"true" required:"true"`
	HostPort    string `split_words:"true"`
}

var (
	Cfg config
)
