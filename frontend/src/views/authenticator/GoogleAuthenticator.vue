<template>
  <div>
    <!-- 上传二维码区域 -->
    <ImportFileDrop
      @qrImported="handleQRImport"
      @cleared="handleClear"
      :acceptJSON="false"
      v-model:has-file="hasFile"
    />

    <!-- URL 输入框 -->
    <n-input
      v-model:value="url"
      placeholder="输入或导入URL"
      clearable
      @blur="onUrlConfirm"
      @keyup.enter="onUrlConfirm"
      :disabled="hasFile"
    />
    <template v-if="data.length!=0">
    <!-- 分类选择 -->
    <n-form-item path="category" label="分类">
      <n-select
        v-model:value="categories"
        filterable
        tag
        :options="categoriesOptions"
      />
    </n-form-item>

    <!-- 表格展示识别结果 -->
    <n-data-table
      v-model:checked-row-keys="selectedRowKeys"
      :columns="columns"
      :data="data"
      :pagination="pagination"
      :row-key="(row: any) => row.id"
      striped
      bordered
      size="small"
      single-select="false"
      @update:checked-row-keys="onSelectRow"
      :loading="loading"
      style="margin-top: 16px"
    />

    <!-- 批量导入按钮 -->
    <n-space justify="end" style="margin-top: 12px">
      <n-button
        type="primary"
        @click="importSelectedRows"
        :disabled="selectedRowKeys.length === 0"
      >
        导入选中项
      </n-button>
    </n-space>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from "vue";
import { useMessage } from "naive-ui";
import ImportFileDrop from "./ImportFileDrop.vue";
import { ParseMigrationURI, ImportTOTP } from "../../../wailsjs/go/app/TOTPHandler";
const message = useMessage();

// ---------- 状态 ----------
const url = ref("");
const lastUrl = ref("");
const data = ref<any[]>([]);
const selectedRowKeys = ref<number[]>([]);
const loading = ref(false);
const categories = ref("未分类");
const hasFile = ref(false);
const props = defineProps<{ categoriesList: any }>();
const categoriesOptions = computed(() =>
  props.categoriesList.map((k: any) => ({ label: k, value: k }))
);
const pagination = { pageSize: 15 };
const emit = defineEmits(["success", "close"]);

// ---------- 输入框逻辑 ----------
function onUrlConfirm() {
  const newUrl = url.value.trim();
  if (newUrl && newUrl !== lastUrl.value) {
    handleQRImport(newUrl);
    lastUrl.value = newUrl;
  }
}

function handleClear() {
  url.value = "";
  lastUrl.value = "";
  data.value = [];
  selectedRowKeys.value = [];
}

// ---------- 核心解析逻辑 ----------
async function handleQRImport(rawData: string) {
  url.value = rawData;
  lastUrl.value = rawData;

  if (!rawData.startsWith("otpauth-migration://")) {
    return message.error("无效的 Google Authenticator 链接");
  }

  try {
    loading.value = true;
    const entries: any[] = await ParseMigrationURI(rawData);

    if (!entries?.length) {
      return message.warning("二维码中未解析到任何账户");
    }

    const now = Date.now();
    entries.forEach((e, i) => (e.id = now + i));

    data.value = entries;
    selectedRowKeys.value = [];

    await nextTick();
    message.success(`成功解析 ${entries.length} 个账户`);
  } catch (err) {
    message.error("无法解析 Google Authenticator 数据");
  } finally {
    loading.value = false;
  }
}

// ---------- 表格逻辑 ----------
function onSelectRow(keys: number[]) {
  selectedRowKeys.value = keys;
}

// ---------- 导入逻辑 ----------
async function importSelectedRows() {
  if (!selectedRowKeys.value.length) {
    return message.warning("请先选择要导入的项目");
  }

  const selected = data.value.filter((item) =>
    selectedRowKeys.value.includes(item.id)
  );
  for (let d of selected) {
    await handleImport(d);
  }
  emit("success");
  message.success("✅ 导入成功");
  emit("close");
}

async function handleImport(form: any) {
  loading.value = true;
  try {
    await ImportTOTP(
      form.account.trim(),
      form.secret.trim(),
      form.remark.trim(),
      form.issuer.trim(),
      categories.value
    );
  } catch (err) {
    message.error(
      "❌ 导入失败：" + (err instanceof Error ? err.message : "未知错误")
    );
  } finally {
    loading.value = false;
  }
}

// ---------- 表格列 ----------
const columns = [
  { type: "selection", multiple: true },
  { title: "颁发者", key: "issuer" },
  { title: "账户", key: "account" },
  { title: "备注", key: "remark" },
];
</script>
