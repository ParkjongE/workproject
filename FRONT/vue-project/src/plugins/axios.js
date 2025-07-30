import axios from 'axios';

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8080',
  timeout: 10000,
  withCredentials: true
});

// access_token이 있으면 Authorization 헤더에 추가
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('jwt_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
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