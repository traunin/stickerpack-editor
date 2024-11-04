import { createRouter, createWebHashHistory } from "vue-router";
import MyStickerpacks from "@/views/MyStickerpacks.vue";

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        { path: "/", redirect: '/my-stickerpacks' },
        { path: "/my-stickerpacks", component: MyStickerpacks },
    ],
});

export default router;