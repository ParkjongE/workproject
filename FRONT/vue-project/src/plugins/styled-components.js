// CSS 클래스를 동적으로 생성하는 함수
export function createStyledComponent(styles, tag = 'div') {
  return {
    name: 'StyledComponent',
    props: {
      className: {
        type: String,
        default: ''
      }
    },
    render() {
      const classes = ['styled-component', this.className].filter(Boolean).join(' ');
      return this.$createElement(tag, {
        class: classes,
        style: styles
      }, this.$slots.default);
    }
  }
}

// styled-components 플러그인
export default {
  install(app) {
    // 전역 메서드로 등록
    app.config.globalProperties.$createStyledComponent = createStyledComponent
  }
} 