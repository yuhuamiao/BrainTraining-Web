<!-- file: src/views/Login.vue -->
<template>
    <div class="login-container">
      <div class="login-form">
        <h2>登录</h2>
        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              v-model="username"
              type="text"
              required
              placeholder="请输入用户名"
            />
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              placeholder="请输入密码"
            />
          </div>
          <button type="submit" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
          <p class="register-link">
            没有账号？<router-link to="/register">立即注册</router-link>
          </p>
        </form>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/authStore'
  
  export default {
    name: 'LoginView',
    setup() {
      const username = ref('')
      const password = ref('')
      const loading = ref(false)
      const error = ref('')
      
      const authStore = useAuthStore()
      const router = useRouter()
      
      const handleLogin = async () => {
        if (!username.value || !password.value) {
          error.value = '请输入用户名和密码'
          return
        }
        
        loading.value = true
        error.value = ''
        
        const result = await authStore.loginUser({
          username: username.value,
          password: password.value
        })
        
        if (result.success) {
          // 登录成功，跳转到首页或之前尝试访问的页面
          const redirect = router.currentRoute.value.query.redirect || '/'
          router.push(redirect)
        } else {
          error.value = result.error
        }
        
        loading.value = false
      }
      
      return {
        username,
        password,
        loading,
        error,
        handleLogin
      }
    }
  }
  </script>
  
  <style scoped>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f5f5;
  }
  
  .login-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }
  
  .login-form h2 {
    text-align: center;
    margin-bottom: 1.5rem;
    color: #333;
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }
  
  .form-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }
  
  button {
    width: 100%;
    padding: 0.75rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 1rem;
  }
  
  button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
  
  .register-link {
    text-align: center;
    margin-top: 1rem;
  }
  
  .error-message {
    color: #f44336;
    margin-top: 1rem;
    text-align: center;
  }
  </style>