<template>
  <n-tabs type="line" animated>
    <n-tab-pane name="import" tab="导入单个">
      <n-form
        ref="formRef"
        :model="form"
        :rules="rules"
        :label-width="100"
        style="text-align: start"
      >
        <n-form-item path="issuer" label="颁发者">
          <n-input v-model:value="form.issuer" placeholder="例如：GitHub" />
        </n-form-item>
        <n-form-item path="account" label="账户">
          <n-input
            v-model:value="form.account"
            placeholder="例如：example@mail.com"
          />
        </n-form-item>
        <n-form-item path="category" label="分类">
          <n-select
            v-model:value="form.category"
            filterable
            tag
            :options="categoriesOptions"
          />
        </n-form-item>
        <n-form-item path="secret" label="密钥">
          <div style="display: flex; flex-direction: column; width: 100%">
            <n-input v-model:value="form.secret" placeholder="输入或导入密钥"  :disabled="hasFile" clearable/>
            <n-collapse v-model:value="showImport" style="margin-bottom: 16px">
              <n-collapse-item name="import">
                <template #header>
                  导入密钥（支持 JSON 或二维码图片）
                </template>
                <ImportFileDrop
                  @jsonImported="handleJSONImport"
                  @qrImported="handleQRImport"
                  :disabled="loading"
                  v-model:has-file="hasFile"
                />
              </n-collapse-item>
            </n-collapse>
          </div>
        </n-form-item>
        <n-form-item path="remark" label="备注">
          <n-input v-model:value="form.remark" placeholder="例如：工作登录" />
        </n-form-item>
      </n-form>

      <n-space justify="end" style="margin-top: 16px">
        <n-button quaternary @click="emit('close')" :disabled="loading"
          >取消</n-button
        >
        <n-button type="primary" @click="handleImport" :loading="loading"
          >导入</n-button
        >
      </n-space>
    </n-tab-pane>
    <n-tab-pane name="GoogleAuthenticator" tab="GoogleAuthenticator 批量导入">
      <GoogleAuthenticator
        @close="emit('close')"
        @success="emit('success')"
        :categories-list="categoriesList"
      ></GoogleAuthenticator>
    </n-tab-pane>
  </n-tabs>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { useMessage } from "naive-ui";
import { ImportTOTP } from "../../../wailsjs/go/app/TOTPHandler";
import ImportFileDrop from "./ImportFileDrop.vue";
import GoogleAuthenticator from "./GoogleAuthenticator.vue";
const emit = defineEmits(["success", "close"]);
const message = useMessage();
const props = defineProps<{ categoriesList: any }>();
const categoriesOptions = computed(() =>
  props.categoriesList.map((k: any) => ({ label: k, value: k }))
);
const hasFile = ref(false);
const formRef = ref();
const loading = ref(false);
const showImport = ref(false);
const form = reactive({
  issuer: "",
  account: "",
  remark: "",
  secret: "",
  category: "未分类",
});

const rules = {
  issuer: { required: true, message: "请输入颁发者", trigger: "blur" },
  account: { required: true, message: "请输入账户", trigger: "blur" },
  secret: { required: true, message: "请输入密钥", trigger: "blur" },
};

function resetForm() {
  form.issuer = "";
  form.account = "";
  form.secret = "";
  form.remark = "";
  form.category = "未分类";
}

function handleJSONImport(entries: any[]) {
  if (entries.length > 0) {
    const e = entries[0];
    form.issuer = e.issuer || "";
    form.account = e.account || "";
    form.secret = e.secret || "";
    form.remark = e.remark || "";
    form.category = e.category || "";
  }
}

function handleQRImport(data: string) {
  if (data.startsWith("otpauth://")) {
    try {
      const url = new URL(data);
      form.account = decodeURIComponent(url.pathname.slice(1));
      const params = new URLSearchParams(url.search);
      form.issuer = params.get("issuer") || "";
      form.secret = params.get("secret") || "";
      message.success("二维码识别成功");
    } catch {
      message.error("无效的 OTPAuth 链接");
    }
  } else {
    message.error("无效的 OTPAuth 链接");
  }
}

async function handleImport() {
  if (loading.value) return;

  await formRef.value?.validate(async (errors: any) => {
    if (errors) return message.warning("请检查表单填写内容");

    loading.value = true;
    try {
      const result = await ImportTOTP(
        form.account.trim(),
        form.secret.trim(),
        form.remark.trim(),
        form.issuer.trim(),
        form.category.trim()
      );
      emit("success", result);
      message.success("✅ 导入成功");
      resetForm();
      emit("close");
    } catch (err) {
      message.error(
        "❌ 导入失败：" + (err instanceof Error ? err.message : "未知错误")
      );
    } finally {
      loading.value = false;
    }
  });
}
</script>
