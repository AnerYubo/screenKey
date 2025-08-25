<template>
  <n-form ref="formRef" :model="form" :rules="rules" :label-width="120" style="text-align: start;">
    <n-form-item label="目标 IP/域名" path="host">
      <n-input v-model:value="form.host" />
    </n-form-item>
    <n-form-item label="目标端口" path="targetPort">
      <n-input-number v-model:value="form.targetPort" :min="1" :max="65535" />
    </n-form-item>
    <n-form-item label="敲门端口序列" path="knockPorts">
      <n-input v-model:value="form.knockPorts" placeholder="如 7000,8000,9000" />
    </n-form-item>
    <n-form-item label="注释" path="remark">
      <n-input v-model:value="form.remark" />
    </n-form-item>
    <n-space justify="end">
      <n-button quaternary @click="emit('close')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">提交</n-button>
    </n-space>
  </n-form>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useMessage } from 'naive-ui';
import type { FormInst } from 'naive-ui';
import { CreateKnock } from '../../../wailsjs/go/app/KnockHandler';

const emit = defineEmits(['success', 'close']);
const props = defineProps<{ rules: any | null }>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);
const form = reactive({
  host: '',
  targetPort: 22,
  knockPorts: '',
  remark: '',
});

const rules = props.rules;

const handleSubmit = async () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      try {
        const item :any= {
          id: Date.now(),
          host: form.host,
          targetPort: form.targetPort,
          knockPorts: form.knockPorts
            .split(',')
            .map((s) => parseInt(s.trim()))
            .filter((v) => !isNaN(v)),
          remark: form.remark,
        };
        // 调用后端接口保存数据
        const err = await CreateKnock(item);
        if(err != null) throw new Error();

        emit('success');
      } catch (err) {
        message.error('保存失败: ' + (err as Error).message);
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