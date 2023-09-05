package config

type SMS struct {
	Limit                  int    `yaml:"limit"`
	WhiteListedPhoneNumber string `yaml:"whitelisted_phone_number"`
}
