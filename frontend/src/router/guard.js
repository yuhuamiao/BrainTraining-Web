import { isTokenValid } from '@/utils/auth'

// 认证守卫 - 保护需要登录的路由
export const authGuard = (to, from, next) => {
  const isAuthenticated = isTokenValid()
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 重定向到登录页，并保存目标路由以便登录后跳转
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else if (to.path === '/login' && isAuthenticated) {
    // 如果已登录但访问登录页，重定向到首页
    next('/')
  } else {
    next()
  }
}

// 角色守卫 - 基于用户角色控制访问（如果需要）
export const roleGuard = (allowedRoles) => {
  return (to, from, next) => {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    
    if (allowedRoles.includes(user.role)) {
      next()
    } else {
      next('/unauthorized')
    }
  }
}