package knock

import (
	"fmt"
	"net"
	"strings"
	"time"

	"screenKey/sqlite"
	"screenKey/utils"
)

type KnockData struct {
	ID         int64  `json:"id"`
	Host       string `json:"host"`
	TargetPort int    `json:"targetPort"`
	KnockPorts []int  `json:"knockPorts"`
	Remark     string `json:"remark"`
}

// 创建新的敲门配置
func CreateKnock(data KnockData) error {
	portsStr := portsToString(data.KnockPorts)
	_, err := sqlite.DB.Insert("knocks", map[string]interface{}{
		"host":        data.Host,
		"target_port": data.TargetPort,
		"knock_ports": portsStr,
		"remark":      data.Remark,
	})
	return err
}

// 获取所有配置并检测端口状态
func ListKnocks() ([]KnockData, error) {
	rows, err := sqlite.DB.GetAll("SELECT * FROM knocks;")
	if err != nil {
		return nil, err
	}

	list := make([]KnockData, len(rows))
	for i, row := range rows {
		list[i] = KnockData{
			ID:         row["id"].(int64),
			Host:       row["host"].(string),
			TargetPort: int(row["target_port"].(int64)),
			KnockPorts: parsePorts(row["knock_ports"].(string)),
			Remark:     row["remark"].(string),
		}
	}

	return list, nil
}

// 删除配置
func DeleteKnock(id int64) error {
	_, err := sqlite.DB.Exec("DELETE FROM knocks WHERE id = ?;", id)
	return err
}

// 更新配置
func UpdateKnock(id int64, data KnockData) error {
	portsStr := portsToString(data.KnockPorts)
	_, err := sqlite.DB.Exec(`UPDATE knocks SET host=?, target_port=?, knock_ports=?, remark=? WHERE id=?`,
		data.Host, data.TargetPort, portsStr, data.Remark, id)
	return err
}

// 敲门逻辑（这里只是打印）
func KnockTarget(data KnockData, proxyAddress string) error {
	fmt.Printf("🔐 开始敲门 %s -> %v -> %d\n", data.Host, data.KnockPorts, data.TargetPort)

	// 先解析域名，拿 IP
	ips, err := net.LookupIP(data.Host)
	if err != nil || len(ips) == 0 {
		return fmt.Errorf("DNS解析失败 %s: %v", data.Host, err)
	}
	ipStr := ips[0].String()
	fmt.Printf("🔍 解析到IP: %s\n", ipStr)
	utils.InitProxy(proxyAddress) // 示例代理地址（Clash/V2Ray 常用）
	// 敲 knockPorts（顺序重要）
	for _, port := range data.KnockPorts {
		address := fmt.Sprintf("%s:%d", ipStr, port)
		conn, err := utils.DialWithTimeout("tcp", address, 1*time.Millisecond) // 超时200ms
		if err != nil {
			fmt.Printf("⚠️ 敲门端口失败 [%s]: %v\n", address, err)
		} else {
			fmt.Printf("✅ 敲门端口成功 [%s]\n", address)
			conn.Close()
		}
		time.Sleep(200 * time.Millisecond) // 每次敲门间隔
	}
	time.Sleep(1 * time.Second) // 敲门结束需要间隔

	// 最终连接目标端口，验证是否开放
	// 最终连接目标端口，验证是否开放
	portOpen := CheckPortOpen(data, proxyAddress)
	if portOpen {
		fmt.Printf("🎉 敲门成功，目标端口 [%s] 已开放！\n", ipStr+":"+fmt.Sprint(data.TargetPort))
		return nil
	} else {
		fmt.Printf("❌ 敲门流程完成，但目标端口 [%s] 并未开放。\n", ipStr+":"+fmt.Sprint(data.TargetPort))
		return fmt.Errorf("目标端口未开放")
	}
}

// 检查端口是否开放
func CheckPortOpen(data KnockData, proxyAddress string) bool {
	address := fmt.Sprintf("%s:%d", data.Host, data.TargetPort)
	fmt.Printf("🚪 正在检测端口是否开放: %v\n", address)

	// 初始化代理（可失败）
	if err := utils.InitProxy(proxyAddress); err != nil {
		fmt.Printf("⚠️ 代理初始化失败: %v\n", err)
		return false
	}

	// 建立 TCP 连接（支持代理）
	conn, err := utils.DialWithTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Printf("❌ 无法连接到 %s: %v\n", address, err)
		return false
	}
	defer conn.Close()

	// 发送 HTTP 请求模拟访问
	request := "GET / HTTP/1.0\r\n\r\n"
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("⚠️ 写入失败，可能端口无服务: %v\n", err)
		return false
	}
	fmt.Printf("✅ 已写入检测请求...\n")

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("⚠️ 无法读取响应，可能端口无服务或协议不兼容: %v\n", err)
		return false
	}

	resp := string(buf[:n])
	fmt.Printf("📥 收到响应（部分内容）:\n%s\n", resp)
	return true
}

// 端口数组转字符串
func portsToString(ports []int) string {
	var parts []string
	for _, p := range ports {
		parts = append(parts, fmt.Sprintf("%d", p))
	}
	return strings.Join(parts, ",")
}

// 字符串转端口数组
func parsePorts(s string) []int {
	var ports []int
	for _, str := range SplitAndTrim(s, ",") {
		var p int
		if _, err := fmt.Sscanf(str, "%d", &p); err == nil {
			ports = append(ports, p)
		}
	}
	return ports
}

// 字符串切割并清洗
func SplitAndTrim(s string, sep string) []string {
	var result []string
	for _, part := range strings.Split(s, sep) {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
