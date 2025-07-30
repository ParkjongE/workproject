// axios import 제거 (임시로 사용하지 않음)

// 네비게이션 가드 - 인증 확인
export const requireAuth = async (to, from, next) => {
  console.log('🔒 라우터 가드 실행:', to.path);
  // 개발 모드에서는 모든 페이지 접근 허용
  console.log('✅ 개발 모드: 모든 페이지 접근 허용');
  next();
};

// 전역 가드 - 모든 라우트에서 실행
export const globalGuard = async (to, from, next) => {
  console.log('🌐 전역 가드 실행:', { to: to.path, from: from.path });
  
  // 로그인 페이지는 항상 접근 가능
  if (to.path === '/') {
    next();
    return;
  }
  
  // 회원가입 페이지는 항상 접근 가능
  if (to.path === '/register') {
    next();
    return;
  }
  
  // 가드 페이지는 항상 접근 가능
  if (to.path === '/auth-guard') {
    next();
    return;
  }
  
  // 개발 모드에서는 모든 페이지 접근 허용
  console.log('✅ 개발 모드: 모든 페이지 접근 허용');
  next();
}; 