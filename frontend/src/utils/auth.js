// file: src/utils/auth.js
// 存储和获取token
export const setToken = (token) => {
    localStorage.setItem('token', token)
  }
  
  export const getToken = () => {
    return localStorage.getItem('token')
  }
  
  export const removeToken = () => {
    localStorage.removeItem('token')
  }
  
  // 存储和获取用户信息
  export const setStoredUser = (user) => {
    localStorage.setItem('user', JSON.stringify(user))
  }
  
  export const getStoredUser = () => {
    const userStr = localStorage.getItem('user')
    return userStr ? JSON.parse(userStr) : null
  }
  
  export const removeStoredUser = () => {
    localStorage.removeItem('user')
  }
  
  // 检查token是否有效
  export const isTokenValid = () => {
    const token = getToken()
    if (!token) return false
    
    try {
      // 简单检查token格式，实际应该解析JWT检查过期时间
      const parts = token.split('.')
      if (parts.length !== 3) return false
      
      // 检查过期时间
      const payload = JSON.parse(atob(parts[1]))
      const exp = payload.exp * 1000 // 转换为毫秒
      return Date.now() < exp
    } catch (error) {
      return false
    }
  }