package config

type config struct {
	AccessTokenSecret  string `split_words:"true" required:"true"`
	AllowApiKey        string `split_words:"true" required:"true"`
	AllowOrigin        string `split_words:"true" required:"true"`
	DatabaseURL        string `split_words:"true" required:"true"`
	HostPort           string `split_words:"true"`
	RediscloudURL      string `split_words:"true"`
	RedisPassword      string `split_words:"true"`
	RefreshTokenSecret string `split_words:"true" required:"true"`
	RedisUsername      string `split_words:"true"`
	MaxRequests        int64  `split_words:"true" required:"true"`
}

var (
	Cfg config
)
