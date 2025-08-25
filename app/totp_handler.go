package app

import (
	"context"
	"screenKey/totp"
)

type TOTPHandler struct {
	ctx context.Context
}

func NewTOTPHandler(ctx context.Context) *TOTPHandler {
	return &TOTPHandler{ctx: ctx}
}

func (h *TOTPHandler) GenerateTOTP(req totp.GenerateTOTPRequest) (*totp.TOTPData, error) {
	return totp.GenerateTOTP(req)
}

func (h *TOTPHandler) ImportTOTP(account, secret, remark, issuer, group string) error {
	return totp.ImportTOTP(account, secret, remark, issuer, group)
}

func (h *TOTPHandler) ListTOTPs() ([]totp.TOTPData, error) {
	data, err := totp.ListTOTPs()
	return data, err
}

func (h *TOTPHandler) GetCurrentCode(secret string) (string, error) {
	return totp.GetCurrentCode(secret)
}

func (h *TOTPHandler) DeleteTOTP(index int64) error {
	return totp.DeleteTOTP(index)
}

func (h *TOTPHandler) UpdateTOTP(index int64, data totp.TOTPData) error {
	return totp.UpdateTOTP(index, data)
}

func (h *TOTPHandler) ParseMigrationURI(uri string) ([]totp.TOTPData, error) {
	return totp.ParseMigrationURI(uri)
}

func (h *TOTPHandler) BuildMigrationURI(account []totp.TOTPData) (string, error) {
	return totp.BuildMigrationURL(account)
}
