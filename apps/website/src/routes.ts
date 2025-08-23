import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { 
      path: '/',
      component: () => import('@/views/stickerpack-overview.vue')
    },
    {
      path: '/create',
      component: () => import('@/views/stickerpack-create.vue')
    },
  ],
})

export default router
