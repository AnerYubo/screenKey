<template>
    <n-button @click="set">点击我-修改</n-button>
    <n-button @click="WindowCenter()">居中</n-button>
    <n-button @click="WindowFullscreen()">全屏窗口</n-button>
    <n-button @click="WindowUnfullscreen()">退出全屏</n-button>
    <n-button @click="WindowMaximise()">最大化</n-button>
    <n-button @click="WindowUnmaximise()">还原最大化</n-button>
    <n-button @click="WindowMinimise()">最小化</n-button>
    <n-button @click="WindowUnminimise()">还原最小化</n-button>
</template>

<script setup>
import { useMessage } from 'naive-ui'
import { NButton, NMessageProvider } from 'naive-ui'
import { WindowSetTitle, WindowCenter, WindowUnfullscreen ,WindowFullscreen,WindowMaximise,WindowUnmaximise, WindowUnminimise, WindowMinimise } from '../../wailsjs/runtime'
const message = useMessage()

async function set() {
WindowSetTitle("魔法启动器 - 老公版")
message.success("修改标题成功")
}
</script>
