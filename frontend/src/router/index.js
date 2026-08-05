import { createRouter, createWebHistory } from 'vue-router'
import { isTokenValid } from '@/utils/auth'

const routes = [
  { path: '/', name: 'Home', component: () => import('@/views/HomePage.vue') },
  { path: '/schulte', name: 'Schulte', component: () => import('@/views/SchulteGrid.vue') },
  { path: '/color-words', name: 'ColorWords', component: () => import('@/views/ColorWords.vue') },
  { path: '/memory', name: 'Memory', component: () => import('@/views/MemoryGame.vue') },
  { path: '/bus-game', name: 'BusGame', component: () => import('@/views/BusNumberGame.vue') },
  { path: '/sudoku', name: 'Sudoku', component: () => import('@/views/SudokuGame.vue') },
  { path: '/progress', name: 'Progress', component: () => import('@/views/ProgressDashboard.vue') },
  { path: '/login', name: 'Login', component: () => import('@/views/Login.vue'), meta: { public: true } },
  { path: '/register', name: 'Register', component: () => import('@/views/Register.vue'), meta: { public: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to) => {
  const authenticated = isTokenValid()
  if (!to.meta.public && !authenticated) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }
  if (to.meta.public && authenticated) {
    return { name: 'Home' }
  }
  return true
})

export default router
