<template>
<n-tabs type="line" animated>
    <n-tab-pane name="qrCodeRef" tab="二维码">
      <div class="showViewModal">
        <div ref="qrCodeRef" class="qr-wrapper">
          <n-qr-code :value="generateOtpAuthUrl(data)" :size="180" />
        </div>
        <n-button
          type="primary"
          size="small"
          style="margin-top: 10px"
          @click="handleDownloadQRCode"
        >
          下载二维码
        </n-button>
      </div>
      </n-tab-pane>    

    <n-tab-pane name="secret" tab="密钥">
      <n-input disabled :value="data?.secret || ''" readonly>
        <template #suffix>
          <n-button size="tiny" type="primary" ghost @click="copysecret"
            >复制</n-button
          >
        </template>
      </n-input>
      </n-tab-pane>    

    <n-tab-pane name="json" tab="JSON">
      <div style="text-align: start; overflow-x: auto" class="code-wrapper">
        <n-scrollbar x-scrollable>
          <div style="width: 100%;padding-bottom: 20px;">
            <n-space vertical :size="16">
              <n-code
                language="javascript"
                :code="JSON.stringify(data, null, 4)"
                style="margin-top: 20px; max-height: 300px; overflow-y: auto"
              />
            </n-space>
          </div>
        </n-scrollbar>
      </div>
      <n-button type="primary" size="small" @click="exportJson"
        >导出 JSON</n-button
      >
      </n-tab-pane>    
</n-tabs>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useMessage } from "naive-ui";

interface Props {
  data: any;
}

const props = defineProps<Props>();

const message = useMessage();
const qrCodeRef = ref<HTMLElement | null>(null);

const data = ref(props.data);

watch(
  () => props.data,
  (newVal) => {
    data.value = newVal;
  }
);

function generateOtpAuthUrl(row: any): string {
  if (!row) return "";
  const label = encodeURIComponent(
    `${row.issuer || "TOTP"}:${row.account || ""}`
  );
  const issuer = encodeURIComponent(row.issuer || "TOTP");
  const secret = encodeURIComponent(row.secret || "");
  return `otpauth://totp/${label}?secret=${secret}&issuer=${issuer}`;
}

function exportJson() {
  if (!data.value) return;
  const json = JSON.stringify(
    {
      issuer: data.value.issuer,
      account: data.value.account,
      secret: data.value.secret,
      remark: data.value.remark,
      otpauth: data.value.otpauth,
    },
    null,
    2
  );
  const blob = new Blob([json], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `totp_${data.value.account}.json`;
  link.click();
  URL.revokeObjectURL(url);
}

const handleDownloadQRCode = () => {
  const canvas = qrCodeRef.value?.querySelector("canvas");
  if (canvas) {
    const url = canvas.toDataURL("image/png");
    const a = document.createElement("a");
    a.href = url;
    a.download = `totp_qrcode_${data.value?.account || "totp"}.png`;
    a.click();
  } else {
    message.warning("二维码未渲染完成");
  }
};

function copysecret() {
  if (!data.value) return;
  navigator.clipboard.writeText(data.value.secret);
  message.success("密钥已复制");
}
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
