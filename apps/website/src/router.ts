import { createRouter, createWebHashHistory } from 'vue-router'
import SharedIcon from '@/assets/icons/shared.svg'
import UserIcon from '@/assets/icons/user.svg'
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    alias: '/user',
    name: 'User packs',
    component: () => import('@/views/user-packs.vue'),
    meta: { inNavbar: true, icon: UserIcon },
  },
  {
    path: '/shared',
    name: 'Shared packs',
    component: () => import('@/views/shared-packs.vue'),
    meta: { inNavbar: true, icon: SharedIcon },
  },
  {
    path: '/create',
    component: () => import('@/views/stickerpack-create.vue'),
  },
  {
    path: '/created/:id',
    name: 'packCreated',
    component: () => import('@/views/stickerpack-created.vue'),
  },
]

export const navbarRoutes = routes.filter((route) => route.meta?.inNavbar)

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
