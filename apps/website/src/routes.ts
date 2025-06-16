import { createRouter, createWebHashHistory } from 'vue-router'
import StickerpackOverview from '@/views/stickerpack-overview.vue'
import StickerpackCreate from './views/stickerpack-create.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: StickerpackOverview },
    { path: '/create', component: StickerpackCreate },
  ],
})

export default router
