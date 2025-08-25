<template>
  <div class="totp-wrapper">
    <n-card size="small" bordered hoverable style="border-radius: 12px" class="totp-card">
      <template #header>
        <FloatingSearch v-model:query="searchQuery">
          <div class="header-title">
            <!-- 左侧按钮组 -->
            <n-button-group>
              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button size="small" @click="showAddModal = true">
                    <template #icon>
                      <n-icon><Cube /></n-icon>
                    </template>
                  </n-button>
                </template>
                <span class="icon-font">生成 TOTP</span>
              </n-popover>

              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button size="small" @click="showImportModal = true">
                    <template #icon>
                      <n-icon><Archive /></n-icon>
                    </template>
                  </n-button>
                </template>
                <span class="icon-font">导入 TOTP</span>
              </n-popover>
            </n-button-group>

            <!-- 右侧过滤 & 操作 -->
            <n-button-group>
              <n-select
                v-model:value="selectedCategory"
                filterable
                size="small"
                :options="[
                  { label: '全部', value: '' },
                  ...categoriesList.map((c) => ({ label: c, value: c })),
                ]"
                placeholder="选择分类过滤"
                style="min-width: 150px; width: auto"
              />

              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    :disabled="selectedRowKeys.length == 0"
                    size="small"
                    @click="selectedRowKeys.length == 1 ? showViewModal = true : showGoogleAuthenticatorViewData = true"
                  >
                    <template #icon>
                      <n-icon><Expand /></n-icon>
                    </template>
                  </n-button>
                </template>
                <span class="icon-font">查看 TOTP</span>
              </n-popover>

              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    :disabled="selectedRowKeys.length !== 1"
                    @click="showEditModal = true"
                  >
                    <template #icon>
                      <n-icon><Create /></n-icon>
                    </template>
                  </n-button>
                </template>
                <span class="icon-font">修改 TOTP</span>
              </n-popover>

              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    :disabled="selectedRowKeys.length === 0"
                    @click="handleDelete"
                  >
                    <template #icon>
                      <n-icon><Trash /></n-icon>
                    </template>
                  </n-button>
                </template>
                <span class="icon-font">删除 TOTP</span>
              </n-popover>
            </n-button-group>
          </div>
        </FloatingSearch>
      </template>

      <!-- ✅ 桌面端表格 -->
      <n-data-table
        v-if="!isMobile"
        v-model:checked-row-keys="selectedRowKeys"
        :columns="columns"
        :data="filteredData"
        :pagination="pagination"
        :row-key="(row: TotpEntry) => row.id"
        striped
        bordered
        size="small"
        single-select="false"
        @update:checked-row-keys="onSelectRow"
        :loading="loading"
        table-layout="auto"
      />

      <!-- ✅ 移动端卡片 -->
      <div v-else class="mobile-cards">
        <TotpCard
          v-for="item in filteredData"
          :key="item.id"
          :data="item"
          :selected="selectedRowKeys.includes(item.id)"
          @update-selection="onCardRightClick"
        />
      </div>
    </n-card>

    <!-- === 各类模态框 === -->
    <myModal v-model:show="showImportModal" title="导入 TOTP 密钥">
      <importData
        v-if="showImportModal"
        @close="() => (showImportModal = false)"
        :categoriesList="categoriesList"
        @success="onAddSuccess"
      />
    </myModal>

    <myModal v-model:show="showAddModal" title="生成 TOTP 密钥">
      <createData
        v-if="showAddModal"
        :categoriesList="categoriesList"
        @close="() => (showAddModal = false)"
        @success="onAddSuccess"
      />
    </myModal>

    <myModal v-model:show="showViewModal" title="查看 TOTP 详情">
      <viewData :data="selectedRow" v-if="showViewModal" />
    </myModal>

    <myModal v-model:show="showEditModal" title="修改 TOTP 密钥">
      <editData
        :data="selectedRow"
        v-if="showEditModal"
        :categoriesList="categoriesList"
        @close="() => (showEditModal = false)"
        @success="onEditSuccess"
      />
    </myModal>
    <myModal v-model:show="showGoogleAuthenticatorViewData" title="批量导出">
      <GoogleAuthenticatorViewData  :data="data.filter((i) => selectedRowKeys.includes(i.id))" v-if="showGoogleAuthenticatorViewData"/>
    </myModal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted, onBeforeUnmount, computed } from "vue";
import FloatingSearch from "@/components/FloatingSearch.vue";
import { useMessage, useDialog } from "naive-ui";
import { Cube, Trash, Expand, Create, Archive } from "@vicons/ionicons5";
import {
  ListTOTPs,
  GetCurrentCode,
  DeleteTOTP,
} from "../../../wailsjs/go/app/TOTPHandler";

import type { TotpEntry } from "./type";
import TotpCard from "./TotpCard.vue";
import editData from "./editData.vue";
import viewData from "./viewData.vue";
import createData from "./createData.vue";
import importData from "./importData.vue";
import myModal from "../../components/myModal.vue";
import CountdownCircle from "@/components/CountdownCircle.vue";
import GoogleAuthenticatorViewData from './GoogleAuthenticatorViewData.vue'
const message = useMessage();
const dialog = useDialog();

const data = ref<TotpEntry[]>([]);
const selectedRowKeys = ref<number[]>([]);
const selectedCategory = ref<string>("");

const selectedRow = computed(() =>
  selectedRowKeys.value.length === 1
    ? data.value.find((i) => i.id === selectedRowKeys.value[0]) || null
    : null
);

function onCardRightClick(id: number) {
  if (selectedRowKeys.value.length === 1 && selectedRowKeys.value[0] === id) {
    selectedRowKeys.value = [];
  } else {
    selectedRowKeys.value = [id];
  }
}

const showAddModal = ref(false);
const showViewModal = ref(false);
const showEditModal = ref(false);
const showImportModal = ref(false);
const showGoogleAuthenticatorViewData = ref(false)
const categoriesList = ref<string[]>([]);

const searchQuery = ref("");

// ✅ 修改 filteredData，加上搜索过滤
const filteredData = computed(() => {
  let result = data.value;
  if (selectedCategory.value) {
    result = result.filter((i) => i.category === selectedCategory.value);
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(
      (i) =>
        i.issuer.toLowerCase().includes(q) ||
        i.account.toLowerCase().includes(q) ||
        i.remark.toLowerCase().includes(q) ||
        i.category.toLowerCase().includes(q)
    );
  }
  return result;
});

const loading = ref(false);
const pagination = { pageSize: 15, showSizePicker: false };

const baseColumns = [
  { type: "selection", multiple: true },
  { title: "颁发者", key: "issuer", ellipsis: { tooltip: true } },
  { title: "账户", key: "account", ellipsis: { tooltip: true } },
  { title: "备注", key: "remark", ellipsis: { tooltip: true } },
  { title: "类别", key: "category", ellipsis: { tooltip: true } },
  {
    title: "验证码",
    align: "right",
    key: "currentCode",
    width: 100,
    render(row: TotpEntry) {
      const isFlashing = row.remainSeconds <= 5;
      return h(
        "div",
        {
          style: {
            display: "flex",
            justifyContent: "flex-end",
            alignItems: "center",
            gap: "8px",
            cursor: "pointer",
            userSelect: "none",
            fontWeight: "bold",
            color: isFlashing ? "red" : "#18a058",
          },
          title: "点击复制验证码",
          onClick: () => copyCurrentCode(row.currentCode),
        },
        [
          h("span", { class: isFlashing ? "flash-red" : "" }, row.currentCode),
          h(CountdownCircle, {
            remainSeconds: row.remainSeconds,
            showNumber: false,
            style: { width: "22px", height: "22px" },
          }),
        ]
      );
    },
  },
];

const columns = computed(() => {
  if (windowWidth.value <= 550) {
    return baseColumns.filter(
      (col) => col.key !== "remark" && col.key !== "category"
    );
  }
  if (windowWidth.value <= 700) {
    return baseColumns.filter((col) => col.key !== "remark");
  }
  return baseColumns;
});

// === 方法区 ===
function copyCurrentCode(code: string) {
  navigator.clipboard.writeText(code);
  message.success("验证码已复制", { duration: 1000 });
}

async function refreshCodes() {
  const now = Date.now();
  for (const item of data.value) {
    item.currentCode = await GetCurrentCode(item.secret);
    item.remainSeconds = 30 - Math.floor((now / 1000) % 30);
  }
}
setInterval(refreshCodes, 1000);

function onSelectRow(keys: number[]) {
  selectedRowKeys.value = keys;
}

async function loadData() {
  loading.value = true;
  try {
    const list = await ListTOTPs();
    const categories: string[] = [];
    if (Array.isArray(list) && list.length > 0) {
      data.value = list.map((item: any) => {
        const category = item.category || "未分类";
        if (!categories.includes(category)) categories.push(category);
        return {
          id: item.id,
          account: item.account,
          issuer: item.issuer,
          secret: item.secret,
          otpauth: item.otpauth,
          remark: item.remark || "",
          category,
          currentCode: "",
          remainSeconds: 0,
        };
      });
      categoriesList.value = categories;
    } else {
      data.value = [];
    }
    await refreshCodes();
  } finally {
    loading.value = false;
  }
}

function onAddSuccess() {
  loadData();
}

function onEditSuccess(updated: TotpEntry) {
  const index = data.value.findIndex((i) => i.id === updated.id);
  if (index !== -1) data.value[index] = { ...data.value[index], ...updated };
  showEditModal.value = false;
  loadData();
}

async function handleDelete() {
  if (loading.value || selectedRowKeys.value.length === 0) return;
  dialog.warning({
    title: () =>
      h("div", { style: "font-weight: bold; font-size: 16px;" }, [
        "确认删除 ",
        h(
          "span",
          { style: "color: red; font-weight: bold" },
          selectedRowKeys.value.length
        ),
        " 条记录？",
      ]),
    content: () =>
      h("div", { style: "color: #555;" }, [
        h("p", null, "这些记录删除后将无法恢复，"),
        h("p", null, "请确认你是否要执行此操作。"),
      ]),
    positiveText: "删除",
    negativeText: "取消",
    async onPositiveClick() {
      loading.value = true;
      try {
        for (const id of selectedRowKeys.value) await DeleteTOTP(id);
        message.success(`✅ 成功删除 ${selectedRowKeys.value.length} 条记录`);
        await loadData();
        selectedRowKeys.value = [];
      } finally {
        loading.value = false;
      }
    },
  });
}

// === 响应式屏幕宽度 ===
const windowWidth = ref(window.innerWidth);
function handleResize() {
  windowWidth.value = window.innerWidth;
}
onMounted(() => {
  loadData();
  window.addEventListener("resize", handleResize);
});
onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
});
const isMobile = computed(() => {
  if (windowWidth.value <= 500) {
    selectedRowKeys.value = [];
    return true;
  }
  return false;
});
</script>

<style scoped lang="scss">
.totp-wrapper {
  padding: 10px;
  max-width: 1000px;
  margin: 0 auto;

  .icon-font {
    font-size: 12px;
    letter-spacing: 2px;
  }
  .totp-card{
    &:hover{
      box-shadow:var(--box-shadow);
    }
  }
}

.header-title {
  display: flex;
  justify-content: space-between;
}

@keyframes flash-red {
  25%,
  100% {
    color: red;
    opacity: 1;
  }
  0%,
  50% {
    color: #ff6666;
    opacity: 0.5;
  }
}
:deep(.flash-red) {
  animation: flash-red 1.5s ease-in-out infinite;
  font-weight: bold;
}

/* 移动端卡片容器 */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
