package handlers

import (
	"myposcore/config"
)

type BaseHandler struct {
	config *config.Config
}

func NewBaseHandler(cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		config: cfg,
	}
}
