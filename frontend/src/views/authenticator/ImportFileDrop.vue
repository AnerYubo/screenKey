<template>
  <n-upload
    v-model:value="fileList"
    :max="1"
    :draggable="true"
    :show-file-list="true"
    :accept="acceptTypes"
    :before-upload="beforeUpload"
    @change="handleChange"
    :directory-dnd="true"
    style="width: 100%; display: flex; flex-direction: column"
  >
    <n-upload-dragger>
      <div
        style="
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          border: 2px dashed var(--n-color-primary);
          cursor: pointer;
          text-align: center;
          padding: 16px;
        "
      >
        <n-icon size="36" color="#0e7a0d">
          <Archive />
        </n-icon>
        <div
          style="
            margin-top: 8px;
            font-size: 14px;
            color: var(--n-text-color-2);
            line-height: 1.4;
          "
        >
          拖拽文件到这里，或点击上传<br />
          仅支持 1 个文件，
          <template v-if="acceptJSON && acceptQR"
            >JSON 文件或二维码图片</template
          >
          <template v-else-if="acceptJSON">JSON 文件</template>
          <template v-else-if="acceptQR">二维码图片</template>
        </div>
      </div>
    </n-upload-dragger>
  </n-upload>
</template>

<script setup lang="ts">
import { ref, computed, defineProps, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { Archive } from "@vicons/ionicons5";
import { readJSONFile, readQRCodeImage } from "./importUtils";

const props = defineProps({
  acceptJSON: {
    type: Boolean,
    default: true,
  },
  acceptQR: {
    type: Boolean,
    default: true,
  },
});

const message = useMessage();
const fileList = ref([]);

const emit = defineEmits<{
  (e: "jsonImported", entries: any[]): void;
  (e: "qrImported", data: string): void;
  (e: "cleared"): void; // 👈 添加 cleared 事件定义
   (e: "update:hasFile", value: boolean): void; // 👈 标准 v-model:has-file
}>();

// 根据 acceptJSON 和 acceptQR 动态设置 accept 属性
const acceptTypes = computed(() => {
  const types: string[] = [];
  if (props.acceptJSON) types.push(".json,application/json");
  if (props.acceptQR) types.push("image/*");
  return types.join(",");
});

function beforeUpload(file: File) {
  const isJson =
    file.type === "application/json" ||
    file.name.toLowerCase().endsWith(".json");
  const isImage = file.type.startsWith("image/");
  if (props.acceptJSON && props.acceptQR) {
    if (!isJson && !isImage) {
      message.error("只支持 JSON 文件或图片");
      return false;
    }
  } else if (props.acceptJSON) {
    if (!isJson) {
      message.error("只支持 JSON 文件");
      return false;
    }
  } else if (props.acceptQR) {
    if (!isImage) {
      message.error("只支持图片");
      return false;
    }
  } else {
    message.error("未配置支持的文件类型");
    return false;
  }

  if (fileList.value.length >= 1) {
    message.warning("只允许上传一个文件");
    return false;
  }
  return true;
}

async function handleChange({
  file,
  fileList,
}: {
  file: any;
  fileList: any[];
}) {
  if (file.status == "removed") {
    emit("cleared");
    updateHasFile(false); // ✅ 移除文件时更新
    return;
  }
  if (!file || file.status == "removed") return; // ✅ 只处理上传完成的

  try {
    const uploadFile = file.file;

    if (!uploadFile) {
      message.error("无法获取上传文件");
      return;
    }
    updateHasFile(true); // ✅ 移除文件时更新

    const isJson =
      uploadFile.type === "application/json" ||
      uploadFile.name.toLowerCase().endsWith(".json");
    const isImage = uploadFile.type.startsWith("image/");

    if (props.acceptJSON && isJson) {
      const entries = await readJSONFile(uploadFile);
      emit("jsonImported", entries);
      message.success(`成功导入 ${entries.length} 条 JSON 记录`);
    } else if (props.acceptQR && isImage) {
      const data = await readQRCodeImage(uploadFile);
      emit("qrImported", data);
    } else {
      message.error("文件类型不支持");
    }
  } catch (err) {
    message.error(
      (err instanceof Error ? err.message : "文件处理失败") || "文件处理失败"
    );
  } finally {
    fileList.splice(0, fileList.length); // 清空上传列表
  }

}
// 封装一个函数来计算是否有文件
function updateHasFile(value: boolean) {
  emit("update:hasFile", value);
}


</script>
