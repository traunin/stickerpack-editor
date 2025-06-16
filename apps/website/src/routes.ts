import { createRouter, createWebHashHistory } from 'vue-router'
import StickerpackOverview from '@/views/stickerpack-overview.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: StickerpackOverview },
  ],
})

export default router
