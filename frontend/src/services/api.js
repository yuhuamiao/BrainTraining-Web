import axios from 'axios'
import { getToken, removeStoredUser, removeToken } from '@/utils/auth'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      removeToken()
      removeStoredUser()
      if (window.location.pathname !== '/login') window.location.assign('/login')
    }
    return Promise.reject(error)
  },
)

export const login = (credentials) => api.post('/login', credentials)
export const register = (userData) => api.post('/register', userData)
export const getMe = () => api.get('/user/me')
export const updateMe = (userData) => api.patch('/user/me', userData)
export const getUserScores = () => api.get('/user/scores')

export const getSchulteMatrix = () => api.get('/schulte/matrix', { params: { size: 5 } })
export const submitSchulteScore = (score) => api.post('/schulte/scores', score)
export const getColorMatrix = (difficulty) => api.get('/color_words/matrix', { params: { difficulty } })
export const submitColorWordScore = (score) => api.post('/color_words/scores', score)
export const getMemoryMatrix = (difficulty) => api.get('/memory/matrix', { params: { difficulty } })
export const submitMemoryScore = (score) => api.post('/memory/scores', score)
export const submitBusScore = (score) => api.post('/bus/scores', score)
export const submitSudokuScore = (score) => api.post('/sudoku/scores', score)

export function apiMessage(error, fallback = '请求失败，请稍后重试') {
  return error.response?.data?.message || error.message || fallback
}

export default api
