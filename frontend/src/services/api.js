// file: src/services/api.js
import axios from 'axios'
import { getToken, removeToken } from '@/utils/auth'

// 创建axios实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加认证token
api.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理认证错误
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token过期或无效，清除本地存储
      removeToken()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证相关API
export const login = (credentials) => {
  return api.post('/login', credentials)
}

export const register = (userData) => {
  return api.post('/register', userData)
}

export const getUserInfo = (userId) => {
  return api.get(`/user/users?userId=${userId}`)
}

export const updateUserInfo = (userData) => {
  return api.post('/user/change', userData)
}

export const uploadAvatar = (formData) => {
  return api.post('/user/photo', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 训练成绩相关API
export const submitSchulteScore = (scoreData) => {
  return api.post('/schulte/scores', scoreData)
}

export const submitColorWordScore = (scoreData) => {
  return api.post('/color_words/scores', scoreData)
}

export const submitMemoryScore = (scoreData) => {
  return api.post('/memory/scores', scoreData)
}

export const submitBusScore = (scoreData) => {
  return api.post('/bus/scores', scoreData)
}

export const getUserScores = (userId) => {
  return api.get(`/user/scores?userId=${userId}`)
}

export default api