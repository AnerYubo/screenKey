package app

import (
	"context"
	"screenKey/knock"
)

type KnockHandler struct {
	ctx context.Context
}

func NewKnockHandler(ctx context.Context) *KnockHandler {
	return &KnockHandler{ctx: ctx}
}

func (h *KnockHandler) CreateKnock(data knock.KnockData) error {
	return knock.CreateKnock(data)
}

func (h *KnockHandler) ListKnocks() ([]knock.KnockData, error) {
	return knock.ListKnocks()
}

func (h *KnockHandler) DeleteKnock(id int64) error {
	return knock.DeleteKnock(id)
}

func (h *KnockHandler) UpdateKnock(id int64, data knock.KnockData) error {
	return knock.UpdateKnock(id, data)
}

func (h *KnockHandler) KnockTarget(data knock.KnockData, proxyAddress string) error {
	return knock.KnockTarget(data, proxyAddress)
}
func (h *KnockHandler) CheckPortOpen(data knock.KnockData, proxyAddress string) bool {
	return knock.CheckPortOpen(data, proxyAddress)
}
