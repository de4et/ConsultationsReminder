package main

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env"`
	StudyType     string `yaml:"studyType"`
	BotKey        string `yaml:"botKey"`
	ChannelChatId int64  `yaml:"channelChatId"`
	GroupNumber   string `yaml:"groupNumber"`
	SchedulePath  string `yaml:"schedulePath"`
}

func ReadConfig(path string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
