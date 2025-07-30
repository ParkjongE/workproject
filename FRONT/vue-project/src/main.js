import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import api from './plugins/axios'
import styledComponentsPlugin from './plugins/styled-components'

// axios 기본 설정
api.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

// 세션 기반 인증을 위한 axios 설정
api.defaults.withCredentials = true; // 쿠키 포함

// 전역 axios 인스턴스 설정
window.axios = api;

// Font Awesome 설정
import { library } from '@fortawesome/fontawesome-svg-core'
import { faSearch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// 사용할 아이콘 추가
library.add(faSearch)

const app = createApp(App)

app.component('font-awesome-icon', FontAwesomeIcon)
app.use(router)
app.use(styledComponentsPlugin)
app.mount('#app')

//Vue 앱의 시작점 — 전역 설정, 라우터, 플러그인, 인터셉터 등록

//axios란? 브라우저나 Node.js에서 HTTP 요청을 보내고 응답을 처리하는 js 라이브러리