package models

import "time"

type Config struct {
	AppSetting        AppSetting    `yaml:"AppSetting"`
	AddrConnection    string        `yaml:"AddrConnection"`
	TimeoutConnection time.Duration `yaml:"TimeoutConnection"`
	RetriesCount      int           `yaml:"RetriesCount"`
}

type AppSetting struct {
	Name   string  `yaml:"Name"`
	Width  float32 `yaml:"Width"`
	Height float32 `yaml:"Height"`
}
