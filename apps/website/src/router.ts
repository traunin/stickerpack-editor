import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
    { 
      path: '/',
      alias: '/stickerpacks',
      name: 'Packs',
      component: () => import('@/views/stickerpack-overview.vue'),
      meta: { inNavbar: true }
    },
    {
      path: '/create',
      component: () => import('@/views/stickerpack-create.vue')
    },
    {
      path: '/created/:id',
      name: "packCreated",
      component: () => import('@/views/stickerpack-created.vue')
    },
]

export const navbarRoutes = routes.filter((route) => route.meta?.inNavbar)

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
