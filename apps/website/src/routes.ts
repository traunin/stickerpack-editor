import { createRouter, createWebHashHistory } from 'vue-router'
import MyStickerpacks from '@/views/MyStickerpacks.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: MyStickerpacks },
  ],
})

export default router
