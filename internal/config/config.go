package config

type Config struct {
	Bot      BotConfig      `map_structure:"bot"`
	Scrapper ScrapperConfig `map_structure:"scrapper"`
}

type BotConfig struct {
	Token string `map_structure:"token"`
	Host  string `map_structure:"host"`
	Port  string `map_structure:"port"`
}

type ScrapperConfig struct {
	BaseURL string `map_structure:"base_url"`
}
