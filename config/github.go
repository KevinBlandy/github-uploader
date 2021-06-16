package config

type GithubConfig struct {
	Owner string			`yaml:"Owner"`
	Repository string		`yaml:"Repository"`
	AccessToken string		`yaml:"AccessToken"`
}