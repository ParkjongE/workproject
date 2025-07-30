import { createRouter, createWebHistory } from 'vue-router'
import api from '@/plugins/axios'

// 도메인별 라우트 import
import authRoutes from './auth.js'
import attendanceRoutes from './attendance.js'
import myScheduleRoutes from './myschedule.js'
import scheduleSharingRoutes from './schedulesharing.js'
import organizationRoutes from './organization.js'

// 모든 라우트 합치기
const routes = [
  { path: '/', redirect: '/login' },
  ...authRoutes,
  ...attendanceRoutes,
  ...myScheduleRoutes,
  ...scheduleSharingRoutes,
  ...organizationRoutes
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// 전역 네비게이션 가드 추가
router.beforeEach(async (to, from, next) => {
  // 로그인/회원가입 페이지는 예외
  if (to.path === '/login' || to.path === '/register') {
    return next();
  }

  // 1. JWT 토큰 체크
  const token = localStorage.getItem('jwt_token');
  if (!token) {
    // 토큰 없으면 로그인 페이지로
    return next({ path: '/login', query: { reason: 'auth' } });
  }

  // 2. 세션 쿠키 체크 (API로 확인)
  try {
    await api.get('/api/user/profile');
    // 세션+JWT 모두 유효하면 통과
    next();
  } catch (e) {
    // 세션 만료/없음 → 로그인 페이지로
    localStorage.removeItem('user');
    localStorage.removeItem('jwt_token');
    next({ path: '/login', query: { reason: 'auth' } });
  }
});

export default router 

//  화면 경로(URL)와 Vue 컴포넌트를 연결                         