<template>
  <div class="auth-layout">
    <section class="auth-brand-panel">
      <div class="auth-brand">
        <span class="brand-mark" aria-hidden="true"><i></i><i></i><i></i><i></i></span>
        <strong>脑力训练台</strong>
      </div>
      <div>
        <p>五项训练</p>
        <h1>把注意力<br />拉回当下</h1>
      </div>
      <span class="auth-footnote">速度 · 记忆 · 心算 · 逻辑</span>
    </section>

    <section class="auth-form-panel">
      <form class="auth-form" @submit.prevent="handleLogin">
        <div>
          <p class="eyebrow">Welcome back</p>
          <h2>登录</h2>
          <p>继续你的训练记录</p>
        </div>
        <div class="field">
          <label for="username">用户名</label>
          <input id="username" v-model.trim="username" class="input" autocomplete="username" required />
        </div>
        <div class="field">
          <label for="password">密码</label>
          <input id="password" v-model="password" class="input" type="password" autocomplete="current-password" required />
        </div>
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <button class="button auth-submit" type="submit" :disabled="loading">{{ loading ? '登录中…' : '登录' }}</button>
        <p class="auth-switch">还没有账号？<router-link to="/register">注册</router-link></p>
      </form>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

async function handleLogin() {
  loading.value = true
  error.value = ''
  const result = await authStore.loginUser({ username: username.value, password: password.value })
  loading.value = false
  if (!result.success) {
    error.value = result.error
    return
  }
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
  router.replace(redirect)
}
</script>

<style src="@/assets/styles/auth.css"></style>
