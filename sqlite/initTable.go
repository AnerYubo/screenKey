package sqlite

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var DB *SQLiteDB

// 获取 Windows 用户本地 AppData 路径
func getUserLocalAppDataDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	localAppData := filepath.Join(usr.HomeDir, "AppData", "Local", "screenKey")
	return localAppData, nil
}

// 初始化数据库，创建表
func InitDatabase() error {
	var err error
	dir, err := getUserLocalAppDataDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	// 创建目录（如果不存在）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	dbPath := filepath.Join(dir, "data.db")

	DB, err = NewSQLiteDB(dbPath)
	if err != nil {
		return err
	}

	// 创建表
	err = DB.CreateTable("totps", map[string]string{
		"id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
		"secret":   "TEXT NOT NULL",
		"otpauth":  "TEXT NOT NULL",
		"issuer":   "TEXT NOT NULL",
		"account":  "TEXT NOT NULL",
		"category": "TEXT",
		"remark":   "TEXT",
	})
	if err != nil {
		return err
	}
	err = DB.CreateTable("knocks", map[string]string{
		"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
		"host":        "TEXT NOT NULL",
		"target_port": "INTEGER NOT NULL",
		"knock_ports": "TEXT NOT NULL",
		"remark":      "TEXT",
	})
	if err != nil {
		return err
	}
	// ✅ 钱包密钥对表
	err = DB.CreateTable("sessions_keys", map[string]string{
		"id":      "INTEGER PRIMARY KEY AUTOINCREMENT",
		"name":    "TEXT NOT NULL",
		"privkey": "TEXT NOT NULL",
		"pubkey":  "TEXT NOT NULL",
		"address": "TEXT NOT NULL",
		"remark":  "TEXT",
	})
	if err != nil {
		return err
	}
	err = DB.CreateTable("contacts", map[string]string{
		"id":      "INTEGER PRIMARY KEY AUTOINCREMENT",
		"name":    "TEXT NOT NULL", // 联系人名称
		"pubkey":  "TEXT NOT NULL", // 公钥（hex 编码）
		"address": "TEXT",          // 可选地址（如果知道）
		"remark":  "TEXT",          // 可选备注
	})
	if err != nil {
		return err
	}
	// ✅ 本地服务实例信息表
	err = DB.CreateTable("servers_instances", map[string]string{
		"id":            "INTEGER PRIMARY KEY AUTOINCREMENT",
		"title":         "TEXT NOT NULL",
		"remark":        "TEXT",
		"privkey_id":    "INTEGER NOT NULL",
		"accessList":    "TEXT",               // JSON字符串
		"auto_sign":     "INTEGER DEFAULT 0",  // 0 表示 false，1 表示 true
		"challenge_ttl": "INTEGER DEFAULT 60", // 签名挑战的有效时间，单位秒
	})

	return err
}
