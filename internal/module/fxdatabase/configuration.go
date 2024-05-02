package fxdatabase

import "github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"

type Configuration struct {
	driver string
	hooks  []hook.Hook
}

func NewConfiguration(driver string, hooks ...hook.Hook) *Configuration {
	return &Configuration{
		driver: driver,
		hooks:  hooks,
	}
}

func (c *Configuration) Driver() string {
	return c.driver
}

func (c *Configuration) Hooks() []hook.Hook {
	return c.hooks
}
