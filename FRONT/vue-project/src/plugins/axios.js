import axios from 'axios';

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8081',
  timeout: 10000,
  withCredentials: true
});

// access_token이 있으면 Authorization 헤더에 추가
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('jwt_token');
    console.log('[Axios Interceptor] 요청 URL:', config.url);
    console.log('[Axios Interceptor] localStorage에서 토큰 확인:', token ? token.substring(0, 20) + '...' : '토큰 없음');
    console.log('[Axios Interceptor] localStorage 전체 키:', Object.keys(localStorage));
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      console.log('[Axios Interceptor] Authorization 헤더 설정 완료');
      console.log('[Axios Interceptor] 설정된 헤더:', config.headers.Authorization.substring(0, 50) + '...');
    } else {
      console.log('[Axios Interceptor] 토큰이 없어 Authorization 헤더 설정 안함');
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 401 응답 처리(재시도 없이 바로 에러 반환)
api.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error)
);

export default api; 