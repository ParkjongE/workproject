<template>
  <InteractionGraphContainer>
    <SideBar />
    <ContentArea>
      <h1>상호작용 그래프</h1>
      <DashboardContainer>
        <GraphView>
          <h2>구성원 간 상호작용 시각화</h2>
          <p>상호작용 그래프 기능은 현재 개발 중입니다.</p>
        </GraphView>
      </DashboardContainer>
    </ContentArea>
  </InteractionGraphContainer>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  InteractionGraphContainer, ContentArea, DashboardContainer, GraphView
} from './InteractionGraph.styled.js';

export default {
  name: 'InteractionGraph',
  components: {
    SideBar,
    InteractionGraphContainer, ContentArea, DashboardContainer, GraphView
  },
  data() {
    return {
      // 데이터 추가 예정
    };
  },
  methods: {
    // 메서드 추가 예정
  }
};
</script> 