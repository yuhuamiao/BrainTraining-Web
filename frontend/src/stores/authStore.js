import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getMe, login, register } from '@/services/api'
import {
  getStoredUser,
  getToken,
  removeStoredUser,
  removeToken,
  setStoredUser,
  setToken,
} from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())
  const user = ref(getStoredUser())
  const isAuthenticated = computed(() => Boolean(token.value))

  function setAuth(authToken, userData) {
    token.value = authToken
    user.value = userData
    setToken(authToken)
    setStoredUser(userData)
  }

  function setUser(userData) {
    user.value = userData
    setStoredUser(userData)
  }

  function clearAuth() {
    token.value = null
    user.value = null
    removeToken()
    removeStoredUser()
  }

  async function loginUser(credentials) {
    try {
      const response = await login(credentials)
      setAuth(response.token, response.user)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.message || '登录失败' }
    }
  }

  async function registerUser(userData) {
    try {
      const response = await register(userData)
      return { success: response.success, message: response.message }
    } catch (error) {
      return { success: false, error: error.response?.data?.message || '注册失败' }
    }
  }

  async function refreshUser() {
    const response = await getMe()
    setUser(response)
    return response
  }

  return { token, user, isAuthenticated, setUser, clearAuth, loginUser, registerUser, refreshUser }
})
