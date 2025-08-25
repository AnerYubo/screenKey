<template>
  <div class="table-wrapper">
    <n-card size="small" bordered hoverable style="border-radius: 12px" class="portknocking-card">
      <template #header>
        <div class="heard-title">
          <n-button-group>
            <n-popover trigger="hover" :delay="200" placement="bottom">
              <template #trigger>
                <n-button size="small" @click="showCreateModal = true">
                  <template #icon
                    ><n-icon><create /></n-icon
                  ></template>
                </n-button>
              </template>
              <span>创建</span>
            </n-popover>
          </n-button-group>
          <div style="display: flex">
            <transition name="slide-fade" mode="out-in">
              <n-input
                v-show="netWork === 'proxy'"
                v-model:value="proxyAddress"
                type="text"
                size="small"
                placeholder="如 127.0.0.1:7890"
                style="width: 160px"
                key="proxy-input"
              />
            </transition>
            <n-button-group>
              <n-button
                size="small"
                :style="{
                  color: netWork == 'proxy' ? '#0e7a0d' : '#afafaf',
                  fontWeight: netWork == 'proxy' ? 'bold' : '',
                }"
                @click="netWork = 'proxy'"
              >
                <span class="icon-font">代理</span>
              </n-button>
              <n-button
                size="small"
                @click="netWork = 'host'"
                :style="{
                  color: netWork == 'host' ? '#0e7a0d' : '#afafaf',
                  fontWeight: netWork == 'host' ? 'bold' : '',
                }"
              >
                <span class="icon-font">直连</span>
              </n-button>
              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    @click="handleKnockAll"
                    :disabled="data.length === 0"
                  >
                    <template #icon
                      ><n-icon color="#0e7a0d"><RocketSharp /></n-icon
                    ></template>
                  </n-button>
                </template>
                <span>一键敲门</span>
              </n-popover>
              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    @click="handleCheckAllPorts"
                    :disabled="data.length === 0"
                  >
                    <template #icon
                      ><n-icon color="#0e7a0d"><Wifi /></n-icon
                    ></template>
                  </n-button>
                </template>
                <span>一键检测信号</span>
              </n-popover>
              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    @click="handleEditOpen"
                    :disabled="!hasSingleSelected"
                  >
                    <template #icon
                      ><n-icon><Create /></n-icon
                    ></template>
                  </n-button>
                </template>
                <span>修改配置</span>
              </n-popover>
              <n-popover trigger="hover" :delay="200" placement="bottom">
                <template #trigger>
                  <n-button
                    size="small"
                    :disabled="!hasSelected"
                    @click="handleDelete"
                  >
                    <template #icon
                      ><n-icon><del_icon /></n-icon
                    ></template>
                  </n-button>
                </template>
                <span>删除配置</span>
              </n-popover>
            </n-button-group>
          </div>
        </div>
      </template>

      <n-data-table
        v-model:checked-row-keys="checkedRowKeys"
        :row-key="(row: RowData) => row.id"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        striped
        bordered
        size="small"
        style="border-radius: 12px"
      />
    </n-card>

    <!-- 创建模态框 -->
    <myModal v-model:show="showCreateModal" title="新增敲门配置">
      <createData
        v-if="showCreateModal"
        :rules="rules"
        @success="onCreateSuccess"
        @close="() => (showCreateModal = false)"
      />
    </myModal>

    <!-- 编辑模态框 -->
    <myModal v-model:show="showEditModal" title="编辑敲门配置">
      <editData
        v-if="showEditModal"
        :model-value="editForm"
        :rules="rules"
        @success="onEditSuccess"
        @close="() => (showEditModal = false)"
      />
    </myModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from "vue";
import {
  LockClosed,
  LockOpen,
  RocketSharp,
  Create,
  Wifi,
} from "@vicons/ionicons5";
import del_icon from "@/components/icons/delete.vue";
import create from "@/components/icons/create.vue";
import type { DataTableColumns, FormRules } from "naive-ui";
import { useMessage, NIcon, NSpin, useDialog } from "naive-ui";

import myModal from "../../components/myModal.vue";
import createData from "./createData.vue";
import editData from "./editData.vue";
import {
  DeleteKnock,
  ListKnocks,
  KnockTarget,
  CheckPortOpen,
} from "../../../wailsjs/go/app//KnockHandler";
import { useDeleteConfirm } from "../../components/DeleteConfirmation";

interface RowData {
  id: number;
  host: string;
  targetPort: number;
  knockPorts: number[];
  remark: string;
  status: string; // "loading" | "open" | "failed" | "closed" | "error" | "未敲门"
}

const message = useMessage();
const dialog = useDialog();
const data = ref<RowData[]>([]);
let idCounter = 0;
const netWork = ref("host"); // 默认选中“主机”
const proxyAddress = ref("127.0.0.1:7890");
const checkedRowKeys = ref<number[]>([]);
const hasSelected = computed(() => checkedRowKeys.value.length > 0);
const hasSingleSelected = computed(() => checkedRowKeys.value.length === 1);

const columns: DataTableColumns<RowData> = [
  { type: "selection", multiple: true },
  { title: "目标 IP/域名", key: "host", align: "center" },
  { title: "目标端口", key: "targetPort", align: "center" },
  {
    title: "状态",
    key: "status",
    align: "center",
    render(row) {
      if (row.status === "loading") {
        return h(
          NSpin,
          { size: "small" },
          { default: () => null } // NSpin 需要一个 slot，哪怕是空的
        );
      }
      const isSuccess = row.status === "open";
      const isFailed =
        row.status === "failed" ||
        row.status === "closed" ||
        row.status === "error" ||
        row.status === "未敲门";
      return h(
        NIcon,
        { color: isSuccess ? "#52c41a" : isFailed ? "#ff4d4f" : undefined },
        () => h(Wifi)
        //() => (isSuccess ? h(LockOpen) : h(LockClosed))
      );
    },
  },
  {
    title: "访问端口",
    key: "knockPorts",
    align: "center",
    render(row) {
      return row.knockPorts.join(", ");
    },
  },
  {
    title: "注释",
    key: "remark",
    align: "center",
    width: 100,
    ellipsis: {
      tooltip: true,
    },
  },
];

const pagination = { pageSize: 15, showSizePicker: false };

const showCreateModal = ref(false);
const showEditModal = ref(false);

const editForm = ref({
  id: 0,
  host: "",
  targetPort: 22,
  knockPorts: "",
  remark: "",
});
const editIndex = ref(-1);

const rules: FormRules = {
  host: [
    { required: true, message: "请输入目标地址", trigger: "blur" },
    {
      validator: (_rule, value) => {
        const domainPattern =
          /^(?!-)[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)*(\.[A-Za-z]{2,})$/;
        const ipPattern =
          /^(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)){3}$/;
        if (!domainPattern.test(value) && !ipPattern.test(value)) {
          return new Error("请输入合法的域名或 IP 地址");
        }
      },
      trigger: "blur",
    },
  ],
  targetPort: {
    required: true,
    validator: (_rule, value) => {
      if (!Number.isInteger(value) || value < 1 || value > 65535) {
        return new Error("端口号必须是 1~65535 的整数");
      }
    },
    trigger: ["input", "blur"],
  },
  knockPorts: {
    required: true,
    validator: (_rule, value) => {
      const ports = value
        .split(",")
        .map((s: string) => parseInt(s.trim()))
        .filter((n: any) => !isNaN(n));
      if (!ports.length || ports.some((p: any) => p < 1 || p > 65535)) {
        return new Error("端口必须为 1~65535 的整数，逗号分隔");
      }
    },
    trigger: "blur",
  },
};

const handleEditOpen = () => {
  if (checkedRowKeys.value.length !== 1) {
    message.warning("请选择一条记录进行编辑");
    return;
  }
  const item = data.value.find((i) => i.id === checkedRowKeys.value[0]);
  if (item) {
    editForm.value = {
      id: item.id,
      host: item.host,
      targetPort: item.targetPort,
      knockPorts: item.knockPorts.join(", "),
      remark: item.remark,
    };
    editIndex.value = data.value.findIndex((i) => i.id === item.id);
    showEditModal.value = true;
  }
};

const handleDelete = async () => {
  if (checkedRowKeys.value.length === 0) {
    message.warning("请先选择要删除的配置");
    return;
  }
  await useDeleteConfirm({
    dialog,
    message,
    count: checkedRowKeys.value.length,
    onConfirm: async () => {
      for (const id of checkedRowKeys.value) {
        const err = await DeleteKnock(id);
        if (err != null) throw new Error(err);
      }
      await loadData();
      checkedRowKeys.value = [];
    },
  });
};

const onCreateSuccess = async () => {
  await loadData();
  showCreateModal.value = false;
  message.success("创建成功");
};

const onEditSuccess = async () => {
  if (editIndex.value !== -1) {
    await loadData();
    showEditModal.value = false;
    editIndex.value = -1;
    message.success("修改成功");
  }
};

// 初始化加载所有敲门配置
const loadData = async () => {
  try {
    const list = await ListKnocks();
    data.value = list.map((item: any) => ({
      ...item,
      status: "未敲门", // 默认状态
      knockPorts: item.knockPorts ?? [],
    }));
    idCounter = data.value.reduce((maxId, item) => Math.max(maxId, item.id), 0);
  } catch (error) {
    message.error("加载敲门配置失败");
  }
};
const types: any = ["success", "info", "warning", "error", "loading"];
// 一键敲门（并发）
// 一键敲门（按 host 分组，同 host 串行，不同 host 并发）
const handleKnockAll = async () => {
  if (data.value.length === 0) {
    message.info("没有敲门配置");
    return;
  }
  const handleKnockAllMsgReactive = message.create("开始一键敲门，请稍候...", {
    type: types[4],
    duration: 1000000,
  });

  // 按 host 分组
  const groupedByHost: Record<string, RowData[]> = {};
  data.value.forEach((item) => {
    if (!groupedByHost[item.host]) {
      groupedByHost[item.host] = [];
    }
    item.status = "loading"; // 初始化为 loading 状态
    groupedByHost[item.host].push(item);
  });

  // 创建每个 host 的串行执行函数
  const taskPromises = Object.values(groupedByHost).map(async (group) => {
    for (const item of group) {
      try {
        const proxy_address = netWork.value == "host" ? "" : proxyAddress.value;
        const err = await KnockTarget(item, proxy_address);
        if (err !== null) throw new Error();
        item.status = "open";
      } catch {
        item.status = "failed";
      }
    }
  });

  // 并发执行各个 host 的串行队列
  await Promise.all(taskPromises);
  handleKnockAllMsgReactive.type = types[0];
  handleKnockAllMsgReactive.content = "一键敲门完成";
  setTimeout(() => {
    handleKnockAllMsgReactive.destroy();
  }, 1000);
};

// 一键检查端口状态（并发）
const handleCheckAllPorts = async () => {
  if (data.value.length === 0) {
    message.info("没有敲门配置");
    return;
  }
  const handleKnockAllMsgReactive = message.create("开始检查端口，请稍候...", {
    type: types[4],
    duration: 10000,
  });
  data.value.forEach((item) => {
    item.status = "loading";
  });
  const proxy_address = netWork.value == "host" ? "" : proxyAddress.value;
  const tasks = data.value.map((item) =>
    CheckPortOpen(item, proxy_address)
      .then((isOpen) => {
        console.log(isOpen);
        item.status = isOpen ? "open" : "closed";
      })
      .catch((err) => {
        item.status = "error";
      })
  );

  await Promise.all(tasks);
  handleKnockAllMsgReactive.type = types[0];
  handleKnockAllMsgReactive.content = "端口检查完成";
  setTimeout(() => {
    handleKnockAllMsgReactive.destroy();
  }, 1000);
};

onMounted(() => {
  loadData();
  //handleKnockAll()
});
</script>

<style scoped lang="scss">
@keyframes spin {
  100% {
    transform: rotate(360deg);
  }
}
.heard-title {
  width: 100%;
  display: flex;
  justify-content: space-between;
  .icon-font {
    font-size: 12px;
    letter-spacing: 2px;
  }
}
.portknocking-card {
  &:hover {
    box-shadow: var(--box-shadow);
  }
}
.table-wrapper {
  padding: 10px;
  max-width: 1000px;
  box-sizing: border-box;
  margin: 0 auto;
}
.n-input-number {
  width: 100%;
}
:deep(.n-spin-body *) {
  width: 16px !important;
  height: 16px !important;
}
/* 动画过渡 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateX(100px);
}
</style>
