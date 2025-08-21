<!-- file: src/views/Register.vue -->
<template>
    <div class="register-container">
      <div class="register-form">
        <h2>注册</h2>
        <form @submit.prevent="handleRegister">
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              v-model="username"
              type="text"
              required
              placeholder="请输入用户名（3-50个字符）"
              minlength="3"
              maxlength="50"
            />
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              placeholder="请输入密码（至少6位）"
              minlength="6"
            />
          </div>
          <div class="form-group">
            <label for="confirmPassword">确认密码</label>
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              required
              placeholder="请再次输入密码"
              minlength="6"
            />
          </div>
          <button type="submit" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>
          <p class="login-link">
            已有账号？<router-link to="/login">立即登录</router-link>
          </p>
        </form>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div v-if="successMessage" class="success-message">
          {{ successMessage }}
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/authStore'
  
  export default {
    name: 'RegisterView',
    setup() {
      const username = ref('')
      const password = ref('')
      const confirmPassword = ref('')
      const loading = ref(false)
      const error = ref('')
      const successMessage = ref('')
      
      const authStore = useAuthStore()
      const router = useRouter()
      
      const handleRegister = async () => {
        // 验证输入
        if (!username.value || !password.value || !confirmPassword.value) {
          error.value = '请填写所有字段'
          return
        }
        
        if (password.value !== confirmPassword.value) {
          error.value = '两次输入的密码不一致'
          return
        }
        
        if (password.value.length < 6) {
          error.value = '密码长度至少6位'
          return
        }
        
        loading.value = true
        error.value = ''
        
        const result = await authStore.registerUser({
          username: username.value,
          password: password.value
        })
        
        if (result.success) {
          successMessage.value = '注册成功，请登录'
          // 2秒后跳转到登录页
          setTimeout(() => {
            router.push('/login')
          }, 2000)
        } else {
          error.value = result.error
        }
        
        loading.value = false
      }
      
      return {
        username,
        password,
        confirmPassword,
        loading,
        error,
        successMessage,
        handleRegister
      }
    }
  }
  </script>
  
  <style scoped>
  .register-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f5f5;
  }
  
  .register-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }
  
  .register-form h2 {
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
    background-color: #2196F3;
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
  
  .login-link {
    text-align: center;
    margin-top: 1rem;
  }
  
  .error-message {
    color: #f44336;
    margin-top: 1rem;
    text-align: center;
  }
  
  .success-message {
    color: #4CAF50;
    margin-top: 1rem;
    text-align: center;
  }
  </style>