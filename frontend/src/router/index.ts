import { createRouter, createWebHistory } from "vue-router";
import { WindowSetTitle } from "../../wailsjs/runtime";
// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "authenticator",
      meta: {
        title: "Authenticator",
      },
      component: () => import("@/views/authenticator/authenticator.vue"), // 异步加载
    },    
    {
      path: "/PortKnocking",
      name: "PortKnocking",
      meta: {
        title: "端口敲门",
      },
      component: () => import("@/views/PortKnocking/PortKnocking.vue"), // 异步加载
    },
    {
      path: "/Setup",
      name: "Setup",
      meta: {
        title: "设置",
      },
      component: () => import("@/views/Setup/Setup.vue"), // 异步加载
    },    
    {
      path: "/About",
      name: "About",
      meta: {
        title: "关于",
      },
      component: () => import("@/views/About/About.vue"), // 异步加载
    },    
  ],    
});

router.beforeEach(async (to, from, next) => {
  const title = "screenKey - 会话钥匙";
  // 解析动态标题
  WindowSetTitle(to.meta.title ? `${to.meta.title} - ${title}` : title);

  next();
});

export default router;
