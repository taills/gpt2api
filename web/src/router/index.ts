import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import BasicLayout from '@/layouts/BasicLayout.vue'
import BlankLayout from '@/layouts/BlankLayout.vue'
import { useUserStore } from '@/stores/user'

/**
 * 路由约定:
 *
 *   meta.public    true  不需要登录
 *   meta.perm      string | string[]  需要任一权限
 *   meta.title     浏览器标签标题
 *
 * 这是前端静态路由表。后端 /api/me/menu 返回的是 UI 菜单,两者并不强绑定:
 * 即便菜单不显示,只要用户输入了正确 URL 且持有权限,也能访问对应页面。
 * 真正的守门人在后端 middleware.RequirePerm,前端只是体验优化。
 */
const routes: RouteRecordRaw[] = [
  // 首页:自带导航 + hero + 界面预览 + 技术栈 + 部署 + 完整 footer,
  // 不套 BlankLayout(避免 BlankLayout 的简短广告 footer 和 Home 自身的 footer 重复)。
  {
    path: '/',
    component: () => import('@/views/landing/Home.vue'),
    meta: { public: true, title: 'GPT2API · ChatGPT 兼容 SaaS 网关 · IMG2 终稿直出 · 批量出图' },
  },
  {
    path: '/',
    component: BlankLayout,
    meta: { public: true },
    children: [
      { path: 'login', component: () => import('@/views/auth/Login.vue'), meta: { public: true, title: '登录' } },
      { path: 'register', component: () => import('@/views/auth/Register.vue'), meta: { public: true, title: '注册' } },
    ],
  },
  {
    path: '/personal',
    component: BasicLayout,
    redirect: '/personal/play',
    children: [
      { path: 'play', component: () => import('@/views/personal/OnlinePlay.vue'),
        meta: { title: '在线体验', perm: ['self:image', 'self:usage'] } },
      { path: 'docs', component: () => import('@/views/personal/ApiDocs.vue'),
        meta: { title: '接口文档', perm: ['self:usage', 'self:image'] } },
      { path: 'keys', component: () => import('@/views/personal/ApiKeys.vue'),
        meta: { title: 'API Keys', perm: 'self:key' } },
      // 旧路径兼容
      { path: 'playground', redirect: '/personal/docs' },
      { path: 'images', redirect: '/personal/play' },
    ],
  },
  {
    path: '/admin',
    component: BasicLayout,
    redirect: '/admin/accounts',
    children: [
      { path: 'accounts', component: () => import('@/views/admin/Accounts.vue'),
        meta: { title: 'GPT账号', perm: 'account:read' } },
      { path: 'image-tasks', component: () => import('@/views/admin/ImageTasks.vue'),
        meta: { title: '生成记录', perm: 'usage:read_all' } },
    ],
  },
  {
    path: '/403',
    component: () => import('@/views/Error403.vue'),
    meta: { public: true, title: '403' },
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/Error404.vue'),
    meta: { public: true, title: '404' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const store = useUserStore()
  const title = (to.meta.title as string) || 'GPT2API 控制台'
  document.title = title

  if (to.meta.public) return true

  if (!store.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  // 还没拉过 me,先补一次(可能来自刷新)
  if (!store.user || store.permissions.length === 0) {
    try {
      await store.fetchMe()
    } catch {
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  }

  const perm = to.meta.perm as string | string[] | undefined
  if (perm && !store.hasPerm(perm)) {
    return { path: '/403' }
  }
  return true
})

export default router
