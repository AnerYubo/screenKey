<!-- src/views/knock/EditData.vue -->
<template>
  <n-form ref="formRef" :model="form" :rules="rules" :label-width="120" style="margin-top: 16px">
    <n-form-item path="host" label="目标 IP/域名">
      <n-input v-model:value="form.host" />
    </n-form-item>

    <n-form-item path="targetPort" label="目标端口">
      <n-input-number v-model:value="form.targetPort" :min="1" :max="65535" />
    </n-form-item>

    <n-form-item path="knockPorts" label="敲门端口序列">
      <n-input v-model:value="form.knockPorts" placeholder="如 7000,8000,9000" />
    </n-form-item>

    <n-form-item path="remark" label="注释">
      <n-input v-model:value="form.remark" />
    </n-form-item>

    <n-space justify="end">
      <n-button quaternary @click="emit('close')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">保存</n-button>
    </n-space>
  </n-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { useMessage } from 'naive-ui';
import type { FormInst } from 'naive-ui';
import {
  UpdateKnock,
} from '../../../wailsjs/go/app//KnockHandler';
import { onMounted } from 'vue';
const emit = defineEmits(['success', 'close']);
const props = defineProps<{
  modelValue: any;
  rules: any;
}>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);
const form = reactive({
  id: 0,
  host: '',
  targetPort: 22,
  knockPorts: '',
  remark: '',
});
onMounted(()=>{
  console.log(props.modelValue)
  Object.assign(form, props.modelValue);
})

const handleSubmit = async () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      try {
        const item = {
          id: form.id,
          host: form.host,
          targetPort: form.targetPort,
          knockPorts: form.knockPorts
            .split(',')
            .map((s) => parseInt(s.trim()))
            .filter((v) => !isNaN(v)),
          remark: form.remark,
        };
        // 调用后端接口
        const err = await UpdateKnock(form.id, item);
        console.log(err)
        if (err != null) {
          throw new Error('更新失败');
        }
        emit('success', item);
      } catch (e) {
        message.error((e as Error).message || '操作失败');
      }
    } else {
      message.error('数据不符合规范');
    }
  });
};

</script>
<style scoped>
.n-input-number {
  width: 100%;
}
</style>