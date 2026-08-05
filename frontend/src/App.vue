<template>
  <div class="app-shell">
    <header v-if="showShell" class="topbar">
      <div class="topbar-inner">
        <router-link class="brand" to="/" aria-label="脑力训练台首页">
          <span class="brand-mark" aria-hidden="true"><i></i><i></i><i></i><i></i></span>
          <span>脑力训练台</span>
        </router-link>

        <nav class="desktop-nav" aria-label="主导航">
          <router-link to="/">训练</router-link>
          <router-link to="/progress">成绩</router-link>
        </nav>

        <div class="account-actions">
          <span class="account-name">{{ authStore.user?.username }}</span>
          <button class="button button-ghost button-small" type="button" @click="logout">退出</button>
        </div>
      </div>
    </header>

    <nav v-if="showShell" class="mobile-nav" aria-label="移动端导航">
      <router-link to="/">训练</router-link>
      <router-link to="/progress">成绩</router-link>
      <button type="button" @click="logout">退出</button>
    </nav>

    <main :class="showShell ? 'app-main' : 'auth-main'">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const showShell = computed(() => !route.meta.public)

function logout() {
  authStore.clearAuth()
  router.replace('/login')
}
</script>
