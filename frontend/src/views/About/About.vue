<template>
  <div class="about-box">
    <n-space vertical size="large">
      <n-card
        title="🧊 关于 screenKey - 会话钥匙"
        size="large"
        bordered
        hoverable
        style="border-radius: 12px"
        class="about-card"
      >
        <n-p>
          <strong>screenKey</strong> 是一个集成了
          <strong>TOTP 动态码、端口敲门</strong>
          的多功能安全工具，旨在提供便捷、高效、极致安全的本地安全能力。
        </n-p>
      </n-card>

      <n-card
        title="📱 TOTP 时间一次性密码"
        bordered
        hoverable
        style="border-radius: 12px"
        class="about-card"

      >
        <n-list bordered>
          <n-list-item>✔ 支持 TOTP 动态码管理</n-list-item>
          <n-list-item>✔ 支持导出 / 导入配置（二维码、密钥、JSON）</n-list-item>
          <n-list-item
            >✔ 支持 Google Authenticator 的批量导入与导出（二维码 / otpauth-migration
            格式）</n-list-item
          >
        </n-list>
      </n-card>

      <n-card
        title="🔐 端口敲门客户端"
        bordered
        hoverable
        style="border-radius: 12px"
        class="about-card"

      >
        <n-list bordered>
          <n-list-item>✔ 可配置服务器地址与敲门端口</n-list-item>
          <n-list-item>✔ 支持多服务配置，一键敲门</n-list-item>
          <n-list-item>
            端口敲门服务端 GitHub 地址：
            <span
              expand-trigger="click"
              line-clamp="3"
              style="color: #0e7a0d; cursor: pointer"
              @click="copyToClipboard('https://github.com/AnerYubo/portknock')"
              title="点击复制地址"
              >https://github.com/AnerYubo/portknock</span
            >
          </n-list-item>
          <n-list-item>
            ✔ 服务端一键安装命令：

              <span
                expand-trigger="click"
                line-clamp="3"
                style="color: #0e7a0d; cursor: pointer"
                title="点击复制命令"
              @click="copyToClipboard(installCmd)"
                >{{ installCmd }}</span
              >
          </n-list-item>
          <n-list-item>✔ 安装简单，开源透明，适合小白无脑部署</n-list-item>
        </n-list>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { Copy } from "@vicons/ionicons5";
import { useMessage } from "naive-ui";
const message = useMessage();
const installCmd = `curl -fsSL https://github.com/aneryubo/portknock//raw//main/install.sh | bash`;
function copyToClipboard(text: any) {
  navigator.clipboard.writeText(text).then(() => {
    message.success("已复制", { duration: 1000 });
  });
}
</script>

<style scoped lang="scss">
.about-box {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}
@media (max-width: 800px) {
  * {
    font-size: 12px;
  }
}
.about-card{
    &:hover{
      box-shadow:var(--box-shadow);
    }
}

</style>
