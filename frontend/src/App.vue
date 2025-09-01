<template>
  <div id="app">
    <!-- 导航栏 -->
    <nav v-if="showNavigation" class="main-nav">
      <div class="nav-brand">大脑训练</div>
      <div class="nav-links">
        <router-link to="/">首页</router-link>
        <router-link to="/bus-game">公交游戏</router-link>
        <!-- 可以添加其他训练链接 -->
      </div>
      <div class="nav-user" v-if="isAuthenticated">
        <span>欢迎, {{ username }}</span>
        <button @click="handleLogout" class="logout-btn">退出</button>
      </div>
      <div class="nav-auth" v-else>
        <router-link to="/login" class="login-btn">登录</router-link>
        <router-link to="/register" class="register-btn">注册</router-link>
      </div>
    </nav>

    <!-- 路由的显示入口 -->
    <main :class="{ 'with-nav': showNavigation }">
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

export default {
  name: 'App',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const authStore = useAuthStore()
    
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const username = computed(() => authStore.user?.username || '')
    
    // 决定是否显示导航栏
    const showNavigation = computed(() => {
      // 在登录和注册页面不显示导航栏
      return route.name !== 'Login' && route.name !== 'Register'
    })
    
    const handleLogout = () => {
      authStore.clearAuth()
      router.push('/login')
    }
    
    return {
      isAuthenticated,
      username,
      showNavigation,
      handleLogout
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: #f8f9fa;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: bold;
  color: #4CAF50;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
}

.nav-links a {
  text-decoration: none;
  color: #2c3e50;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  background-color: #e9ecef;
  color: #4CAF50;
}

.nav-user, .nav-auth {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logout-btn, .login-btn, .register-btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.9rem;
}

.logout-btn {
  background-color: #f44336;
  color: white;
}

.login-btn {
  background-color: #2196F3;
  color: white;
}

.register-btn {
  background-color: #4CAF50;
  color: white;
}

main {
  flex: 1;
  padding: 2rem;
}

main.with-nav {
  padding-top: 1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-nav {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }
  
  .nav-links {
    order: 3;
    width: 100%;
    justify-content: center;
  }
  
  main {
    padding: 1rem;
  }
}
</style>