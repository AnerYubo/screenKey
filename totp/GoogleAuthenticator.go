package totp

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/protobuf/proto"

	pb "screenKey/totp/googleauth" // 这里改成你生成 proto 的 go_package
)

// normalizeSecret 尝试把各种 secret 转成合法的 Base32 bytes
func normalizeSecret(secret string) ([]byte, error) {
	s := strings.TrimSpace(secret)

	// 统一大写（Base32 要求大写 A-Z2-7）
	s = strings.ToUpper(s)

	// 尝试 Base32 解码（无 padding）
	if raw, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s); err == nil {
		return raw, nil
	}

	// 尝试 Base32 解码（有 padding）
	if raw, err := base32.StdEncoding.DecodeString(s); err == nil {
		return raw, nil
	}

	// 尝试 Base64
	if raw, err := base64.StdEncoding.DecodeString(secret); err == nil {
		return raw, nil
	}

	// 尝试 hex
	if raw, err := hex.DecodeString(secret); err == nil {
		return raw, nil
	}

	return nil, fmt.Errorf("unsupported secret format: %s", secret)
}

// BuildMigrationURL 根据多个账号信息生成 otpauth-migration:// URL
func BuildMigrationURL(accounts []TOTPData) (string, error) {
	var otpParams []*pb.MigrationPayload_OTPParameters

	for _, acc := range accounts {
		// 自动修正 secret
		secret, err := normalizeSecret(acc.Secret)
		if err != nil {
			return "", fmt.Errorf("invalid secret %s: %v", acc.Secret, err)
		}

		otpParams = append(otpParams, &pb.MigrationPayload_OTPParameters{
			Secret:    secret,
			Name:      acc.AccountName,
			Issuer:    acc.Issuer,
			Algorithm: pb.MigrationPayload_SHA1,
			Digits:    pb.MigrationPayload_SIX,
			Type:      pb.MigrationPayload_TOTP,
		})
	}

	payload := &pb.MigrationPayload{
		OtpParameters: otpParams,
		Version:       1,
		BatchSize:     int32(len(otpParams)),
		BatchIndex:    0,
		BatchId:       123456, // 可以随机
	}

	data, err := proto.Marshal(payload)
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	escaped := url.QueryEscape(b64)

	return "otpauth-migration://offline?data=" + escaped, nil
}

// ParseMigrationURL 解析 otpauth-migration:// URL，返回 OTPParameters
func ParseMigrationURL(uri string) ([]*pb.MigrationPayload_OTPParameters, error) {
	if !strings.HasPrefix(uri, "otpauth-migration://offline?data=") {
		return nil, fmt.Errorf("invalid otpauth-migration url")
	}

	// 提取 data 参数
	raw := strings.TrimPrefix(uri, "otpauth-migration://offline?data=")

	// URL decode
	decoded, err := url.QueryUnescape(raw)
	if err != nil {
		return nil, fmt.Errorf("url decode failed: %v", err)
	}

	// Base64 decode
	data, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %v", err)
	}

	// Protobuf 反序列化
	payload := &pb.MigrationPayload{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, fmt.Errorf("proto unmarshal failed: %v", err)
	}

	return payload.OtpParameters, nil
}
