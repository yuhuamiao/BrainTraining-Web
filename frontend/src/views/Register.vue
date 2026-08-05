<template>
  <div class="auth-layout">
    <section class="auth-brand-panel register-panel">
      <div class="auth-brand">
        <span class="brand-mark" aria-hidden="true"><i></i><i></i><i></i><i></i></span>
        <strong>脑力训练台</strong>
      </div>
      <div>
        <p>建立基线</p>
        <h1>从今天<br />开始记录</h1>
      </div>
      <span class="auth-footnote">每次练习都会留下进步轨迹</span>
    </section>

    <section class="auth-form-panel">
      <form class="auth-form" @submit.prevent="handleRegister">
        <div>
          <p class="eyebrow">Create account</p>
          <h2>注册</h2>
          <p>创建你的训练档案</p>
        </div>
        <div class="field">
          <label for="username">用户名</label>
          <input id="username" v-model.trim="username" class="input" minlength="3" maxlength="50" autocomplete="username" required />
        </div>
        <div class="field">
          <label for="password">密码</label>
          <input id="password" v-model="password" class="input" type="password" minlength="6" autocomplete="new-password" required />
        </div>
        <div class="field">
          <label for="confirm-password">确认密码</label>
          <input id="confirm-password" v-model="confirmPassword" class="input" type="password" minlength="6" autocomplete="new-password" required />
        </div>
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <p v-if="success" class="form-success">{{ success }}</p>
        <button class="button auth-submit" type="submit" :disabled="loading || Boolean(success)">{{ loading ? '注册中…' : '注册' }}</button>
        <p class="auth-switch">已有账号？<router-link to="/login">登录</router-link></p>
      </form>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')
const router = useRouter()
const authStore = useAuthStore()

async function handleRegister() {
  error.value = ''
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }
  loading.value = true
  const result = await authStore.registerUser({ username: username.value, password: password.value })
  loading.value = false
  if (!result.success) {
    error.value = result.error
    return
  }
  success.value = '注册成功，正在前往登录页…'
  window.setTimeout(() => router.replace('/login'), 800)
}
</script>

<style src="@/assets/styles/auth.css"></style>
