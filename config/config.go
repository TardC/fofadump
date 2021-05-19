package config

type Config struct {
	FofaServer string
	Email      string `json:"email,omitempty"`
	Key        string `json:"key,omitempty"`
}
