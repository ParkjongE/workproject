import Login from '../views/Login/Login.vue'

const authRoutes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/Login/Register.vue')
  },
  {
    path: '/auth-guard',
    name: 'authGuard',
    component: () => import('../views/Login/AuthGuard.vue')
  }
]

export default authRoutes 