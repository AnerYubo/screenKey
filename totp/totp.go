package totp

import (
	"encoding/base32"
	"fmt"
	"log"
	"screenKey/sqlite" // 你的 sqlite 包
	"time"

	"github.com/pquerna/otp/totp"
)

// 请求结构体
type GenerateTOTPRequest struct {
	Issuer      string `json:"issuer"`
	AccountName string `json:"account"`
	Remark      string `json:"remark"`
	Category    string `json:"category"` // 修改字段名：分组 -> 类别
}

// 存储结构体
type TOTPData struct {
	ID          int64  `json:"id"`
	Secret      string `json:"secret"`
	Otpauth     string `json:"otpauth"`
	Issuer      string `json:"issuer"`
	AccountName string `json:"account"`
	Remark      string `json:"remark"`
	Category    string `json:"category"` // 修改字段名：分组 -> 类别
}

// 生成新的 TOTP 密钥
func GenerateTOTP(req GenerateTOTPRequest) (*TOTPData, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      req.Issuer,
		AccountName: req.AccountName,
	})
	if err != nil {
		return nil, fmt.Errorf("生成 TOTP key 失败: %w", err)
	}

	data := &TOTPData{
		Secret:      key.Secret(),
		Otpauth:     key.URL(),
		Issuer:      req.Issuer,
		AccountName: req.AccountName,
		Remark:      req.Remark,
		Category:    req.Category, // 使用传入的类别
	}

	_, err = sqlite.DB.Insert("totps", map[string]interface{}{
		"secret":   data.Secret,
		"otpauth":  data.Otpauth,
		"issuer":   data.Issuer,
		"account":  data.AccountName,
		"remark":   data.Remark,
		"category": data.Category, // 插入类别字段
	})
	if err != nil {
		return nil, fmt.Errorf("插入数据库失败: %w", err)
	}

	log.Printf("生成新 TOTP：Issuer=%s, Account=%s\n", data.Issuer, data.AccountName)
	return data, nil
}

// 导入已有 TOTP
func ImportTOTP(account, secret, remark, issuer, category string) error {
	// 校验密钥有效性
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return fmt.Errorf("密钥无效: %w", err)
	}
	log.Printf("导入 TOTP，当前验证码为：%s\n", code)

	_, err = sqlite.DB.Insert("totps", map[string]interface{}{
		"secret":   secret,
		"otpauth":  fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, account, secret, issuer),
		"issuer":   issuer,
		"account":  account,
		"remark":   remark,
		"category": category, // 添加类别字段
	})
	if err != nil {
		return fmt.Errorf("插入数据库失败: %w", err)
	}

	log.Printf("导入 TOTP 成功: Issuer=%s, Account=%s\n", issuer, account)
	return nil
}

// 列出所有 TOTP
func ListTOTPs() ([]TOTPData, error) {
	rows, err := sqlite.DB.GetAll("SELECT * FROM totps;")
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	var totps []TOTPData
	for _, row := range rows {
		// 类型断言，注意这里可能导致 panic，建议根据你的业务逻辑自行校验
		id, ok := row["id"].(int64)
		if !ok {
			log.Printf("警告: id 字段类型异常: %v\n", row["id"])
			continue
		}
		secret, _ := row["secret"].(string)
		otpauth, _ := row["otpauth"].(string)
		issuer, _ := row["issuer"].(string)
		account, _ := row["account"].(string)
		remark, _ := row["remark"].(string)
		category, _ := row["category"].(string) // 获取类别字段

		totp := TOTPData{
			ID:          id,
			Secret:      secret,
			Otpauth:     otpauth,
			Issuer:      issuer,
			AccountName: account,
			Remark:      remark,
			Category:    category, // 填充分组字段
		}
		totps = append(totps, totp)
	}

	log.Printf("查询到 %d 条 TOTP 记录\n", len(totps))
	return totps, nil
}

// 删除某一项
func DeleteTOTP(id int64) error {
	affected, err := sqlite.DB.Exec("DELETE FROM totps WHERE id = ?;", id)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	log.Printf("删除 TOTP id=%d，影响行数: %d\n", id, affected)
	if affected == 0 {
		log.Printf("警告：未删除任何记录，id %d 可能不存在\n", id)
	}

	return nil
}

// 修改某一项
func UpdateTOTP(id int64, newData TOTPData) error {
	affected, err := sqlite.DB.Exec("UPDATE totps SET secret = ?, otpauth = ?, issuer = ?, account = ?, remark = ?, `category` = ? WHERE id = ?;",
		newData.Secret, newData.Otpauth, newData.Issuer, newData.AccountName, newData.Remark, newData.Category, id) // 更新 `category`
	if err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}
	log.Printf("更新 TOTP id=%d，影响行数: %d\n", id, affected)
	if affected == 0 {
		log.Printf("警告：未更新任何记录，id %d 可能不存在\n", id)
	}
	return nil
}

// 获取当前动态密码
func GetCurrentCode(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("生成验证码失败: %w", err)
	}
	return code, nil
}

// 解析 otpauth-migration URI 并返回 TOTP 列表（不入库）
// ParseMigrationURI 解析 otpauth-migration URI 并返回 TOTPData 列表（不入库）
func ParseMigrationURI(uri string) ([]TOTPData, error) {
	// 调用我们自己写的 ParseMigrationURL
	entries, err := ParseMigrationURL(uri)
	if err != nil {
		return nil, fmt.Errorf("解析 otpauth-migration URI 失败: %w", err)
	}

	var results []TOTPData
	for _, entry := range entries {
		// 转回 Base32 secret，Google Authenticator 常用 Base32
		secretB32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(entry.Secret)

		// 验证 secret 合法性（可选）
		if _, err := totp.GenerateCode(secretB32, time.Now()); err != nil {
			log.Printf("跳过无效条目（%s）: %v\n", entry.Name, err)
			continue
		}

		// 拼 otpauth:// URL（可选）
		otpauthURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d",
			entry.Issuer, entry.Name, secretB32, entry.Issuer, "SHA1", 6,
		)

		totpData := TOTPData{
			ID:          0, // 尚未入库
			Secret:      secretB32,
			Otpauth:     otpauthURL,
			Issuer:      entry.Issuer,
			AccountName: entry.Name,
			Remark:      "",
			Category:    "", // 无分组
		}
		results = append(results, totpData)
	}

	log.Printf("解析 otpauth-migration 成功，共 %d 项\n", len(results))
	return results, nil
}
