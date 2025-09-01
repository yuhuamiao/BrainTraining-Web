import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, getUserInfo } from '@/services/api'
import { setToken, removeToken, getToken } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())
  const user = ref(null)
  const isAuthenticated = computed(() => !!token.value)

  const setUser = (userData) => {
    user.value = userData
  }

  const setAuthToken = (newToken) => {
    token.value = newToken
    setToken(newToken)
  }

  const clearAuth = () => {
    token.value = null
    user.value = null
    removeToken()
  }

  const loginUser = async (credentials) => {
    try {
      const response = await login(credentials)
      if (response.success) {
        setAuthToken(response.token)
        setUser(response.user)
        return { success: true }
      }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || '登录失败' 
      }
    }
  }

  const registerUser = async (userData) => {
    try {
      const response = await register(userData)
      return { success: response.success, message: response.message }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || '注册失败' 
      }
    }
  }

  const fetchUserInfo = async (userId) => {
    try {
      const userInfo = await getUserInfo(userId)
      setUser(userInfo)
      return userInfo
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果token无效，清除认证状态
      if (error.response?.status === 401) {
        clearAuth()
      }
      throw error
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    setUser,
    setAuthToken,
    clearAuth,
    loginUser,
    registerUser,
    fetchUserInfo
  }
})