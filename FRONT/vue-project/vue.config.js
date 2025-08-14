const { defineConfig } = require('@vue/cli-service')
const path = require('path');

module.exports = {
  devServer: {
    host: '0.0.0.0',  // 모든 네트워크 인터페이스에서 접속 허용
    port: 3000,
    allowedHosts: 'all'  // 모든 호스트에서 접속 허용
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@styled': path.join(__dirname, 'src', 'views', 'myschedule', 'Schedule', 'styled.js')
      }
    }
  },
  transpileDependencies: true,
  chainWebpack: config => {
    // HTML 파일을 raw 문자열로 가져오기 위한 설정
    config.module
      .rule('raw')
      .test(/\.html$/)
      .resourceQuery(/raw/)
      .use('raw-loader')
      .loader('raw-loader')
      .end();
  },
  // 백엔드 API 프록시 설정
  devServer: {
    host: '0.0.0.0',  
    port: 3000, // 포트를 3000로 변경
    allowedHosts: 'all',
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true
      }
    }
  }
}


//Vue CLI 프로젝트의 빌드, 개발 서버, 웹팩 등 세부 설정 파일