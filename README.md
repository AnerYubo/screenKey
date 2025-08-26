# 🧊 screenKey - 会话钥匙

一个基于 **Golang + Wails + Vue3 + Vite + Naive UI** 开发的跨平台桌面应用。  
集成了 **TOTP 动态口令管理器** 和 **端口敲门客户端**，为开发者和安全爱好者提供便捷、高效、安全的工具体验。

---

## ✨ 功能特性

### 📱 TOTP 时间一次性密码
- ✔ 支持 TOTP 动态码管理  
- ✔ 支持 **导出 / 导入配置**（二维码、密钥、JSON）  
- ✔ 兼容 **Google Authenticator**，支持批量导入与导出（二维码 / `otpauth-migration` 格式）  
- ✔ 支持 Ctrl + k [ Ctrl + / ] 可唤醒搜索栏

### 🔐 端口敲门客户端
- ✔ 可配置服务器地址与敲门端口  
- ✔ 支持多服务配置，一键敲门  
- ✔ 内置端口检测功能，方便确认服务是否可用  
- ✔ 一键复制 **服务端 GitHub 地址**：[AnerYubo/portknock](https://github.com/AnerYubo/portknock)  
- ✔ 服务端一键安装命令：
  ```bash
  curl -fsSL https://github.com/aneryubo/portknock/raw/main/install.sh | bash

* ✔ 安装简单，开源透明，适合小白无脑部署

### 🎨 界面与主题

* ✔ 支持 **亮色 / 暗色 / 跟随系统** 三种主题模式
* ✔ 自动监听系统主题切换
* ✔ CSS 变量驱动，支持全局 UI 自适应
* ✔ 响应式布局，桌面端表格 + 移动端卡片

---

## 🚀 开始使用

### 克隆项目

```bash
git clone https://github.com/AnerYubo/screenKey.git
cd screenKey
```

### 安装依赖

```bash
go mod tidy
```

### 启动开发环境

```bash
wails dev
```

### 构建生产版本

```bash
wails build
```

---

## 🛠 技术栈

* **后端**：Golang
* **桌面框架**：Wails
* **前端框架**：Vue 3 + Vite
* **UI 组件库**：Naive UI
* **代码高亮**：highlight.js

---

## ⭐ 支持项目

如果你觉得这个工具对你有帮助，欢迎点个 **Star** 支持一下！



