<template>
  <span>支持 Google Authenticator 导入</span>
  <n-tabs type="line" animated v-model:value="tabsValue">
    <n-tab-pane name="qrCodeRef" tab="二维码">
      <div class="showViewModal">
        <div ref="qrCodeRef" class="qr-wrapper">
          <n-qr-code :value="url[page - 1]" :size="180" />
        </div>
        <div style="display: flex;gap: 5px;justify-content: center;">
          <n-button
            type="primary"
            size="small"
            style="margin-top: 10px"
            @click="handleDownloadQRCode"
          >
            下载二维码
          </n-button>
          <n-button
            type="tertiary"
            size="small"
            style="margin-top: 10px"
            @click="validateQRCode"
          >
            检测二维码
          </n-button>
        </div>
      </div>
    </n-tab-pane>

    <n-tab-pane name="url" tab="url">
      <n-input disabled :value="url[page-1] || ''" readonly>
        <template #suffix>
          <n-button size="tiny" type="primary" ghost @click="copysecret"
            >复制</n-button
          >
        </template>
      </n-input>
    </n-tab-pane>
  </n-tabs>
  <div style="margin-top: 20px; display: flex; justify-content: end">
    <n-pagination v-model:page="page" :page-count="url.length" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useMessage } from "naive-ui";
import { BuildMigrationURI } from "../../../wailsjs/go/app/TOTPHandler";
import { getConfig } from "../Setup/init";
import { readQRCodeImage } from "./importUtils";
import { ParseMigrationURI, ImportTOTP } from "../../../wailsjs/go/app/TOTPHandler";

interface Props {
  data: any;
}
const tabsValue = ref("qrCodeRef");
const props = defineProps<Props>();
const url = ref<string[]>([]);
const page = ref(1); // 当前页码
const message = useMessage();
const qrCodeRef = ref<HTMLElement | null>(null);

const data = ref(props.data);

watch(
  () => props.data,
  (newVal) => {
    data.value = newVal;
  }
);

// ---------- 验证二维码是否可识别 ----------
async function validateQRCode() {
  const canvas = qrCodeRef.value?.querySelector("canvas");
  if (!canvas) {
    return message.warning("二维码未渲染完成");
  }

  try {
    // 1. 转成 DataURL
    const dataUrl = canvas.toDataURL("image/png");

    // 2. 转成 File 对象再调用 readQRCodeImage
    const res = await fetch(dataUrl);
    const blob = await res.blob();
    const file = new File([blob], "qrcode.png", { type: "image/png" });

    // 3. 调用识别函数
    const result = await readQRCodeImage(file);
    await handleQRImport(result);
  } catch (err) {
    console.error(err);
    message.error("二维码识别失败;请减少一次性导出数据");
  }
}
// ---------- 核心解析逻辑 ----------
async function handleQRImport(rawData: string) {

  if (!rawData.startsWith("otpauth-migration://")) {
    message.error("无效的 Google Authenticator 链接");
    return new Error("无效的 Google Authenticator 链接");
  }

  try {
    const entries: any[] = await ParseMigrationURI(rawData);

    if (!entries?.length) {
      message.warning("二维码中未解析到任何账户");
      return new Error("无效的 Google Authenticator 链接");

    }

    message.success(`成功解析 ${entries.length} 个账户`);
  } catch (err) {
    message.error("无法解析 Google Authenticator 数据");
    return new Error("无效的 Google Authenticator 链接");

  }
}
const handleDownloadQRCode = () => {
  const canvas = qrCodeRef.value?.querySelector("canvas");
  if (canvas) {
    const url = canvas.toDataURL("image/png");
    const a = document.createElement("a");
    a.href = url;
    a.download = `otpauth_migration_${data.value?.account || "totp"}.png`;
    a.click();
  } else {
    message.warning("二维码未渲染完成");
  }
};

function copysecret() {
  if (!data.value) return;
  navigator.clipboard.writeText(data.value.secret);
  message.success("url已复制");
}
function chunkArray<T>(array: any[], size: number): any[][] {
  const result: T[][] = [];
  for (let i = 0; i < array.length; i += size) {
    result.push(array.slice(i, i + size));
  }
  return result;
}
onMounted(async () => {
  const config = getConfig("Authenticator");

  //const data = await BuildMigrationURI(props.data)
  const chunks = chunkArray(props.data, config.exportCount);

  for (const chunk of chunks) {
    const data = await BuildMigrationURI(chunk);
    url.value.push(data);
  }
});
</script>
<style lang="css" scoped>
.qr-wrapper {
  padding: 20px;
  margin: 0 auto;
  max-width: max-content;
  padding-right: 40px;
  margin-bottom: 12px; /* ✅ 增加底部间距，避免被按钮顶住 */
  background-color: #fff;
  border: 1px solid #eee;
  border-radius: 12px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
}
.code-wrapper {
  padding: 20px;
  margin: 0 auto;
  max-width: max-content;
  margin-bottom: 12px; /* ✅ 增加底部间距，避免被按钮顶住 */
  background-color: #fff;
  border: 1px solid #eee;
  border-radius: 12px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
}
</style>
