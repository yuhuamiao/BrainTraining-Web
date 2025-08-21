// file: src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { authGuard } from './guards'
import HomePage from '../views/HomePage.vue'
import BusNumberGame from '../views/BusNumberGame.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'

const routes = [
  {
    path: '/',
    component: HomePage,
    meta: { requiresAuth: true } // 首页需要认证
  },
  // 公交游戏路由
  {
    path: '/bus-game',
    name: 'BusGame',
    component: BusNumberGame,
    meta: { requiresAuth: true } // 游戏页面需要认证
  },
  // 登录页面
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false } // 登录页面不需要认证
  },
  // 注册页面
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false } // 注册页面不需要认证
  },
  // 可以继续添加其他训练的路由...
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 添加路由守卫
router.beforeEach(authGuard)

export default router