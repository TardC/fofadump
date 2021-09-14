package fofadump

type Config struct {
	FofaServer string
	Email      string `json:"email,omitempty"`
	Key        string `json:"key,omitempty"`
}

func NewFofaConfig() *Config {
	return &Config{
		FofaServer: "https://fofa.so",
		Email:      "",
		Key:        "",
	}
}
