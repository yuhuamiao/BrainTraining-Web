const TOKEN_KEY = 'brain-training-token'
const USER_KEY = 'brain-training-user'

export const setToken = (token) => localStorage.setItem(TOKEN_KEY, token)
export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const removeToken = () => localStorage.removeItem(TOKEN_KEY)
export const setStoredUser = (user) => localStorage.setItem(USER_KEY, JSON.stringify(user))

export function getStoredUser() {
  try {
    return JSON.parse(localStorage.getItem(USER_KEY))
  } catch {
    return null
  }
}

export const removeStoredUser = () => localStorage.removeItem(USER_KEY)

export function isTokenValid() {
  const token = getToken()
  if (!token) return false
  try {
    const encoded = token.split('.')[1]
    if (!encoded) return false
    const normalized = encoded.replace(/-/g, '+').replace(/_/g, '/')
    const payload = JSON.parse(atob(normalized))
    return Number(payload.exp) * 1000 > Date.now()
  } catch {
    return false
  }
}
