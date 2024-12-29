import Signin from '@/views/Signin.vue'
import Chat from '@/views/Chat.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: Signin },
    { path: '/chat', component: Chat , meta: {requiresAuth: true }},
  ],
})

export default router
