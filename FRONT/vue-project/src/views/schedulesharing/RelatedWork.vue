<template>
  <RelatedWorkContainer>
    <SideBar />
    <ContentArea>
      <h1>연관업무 기록 조회</h1>
      <DashboardContainer>
        <WorkRecords>
          <h2>업무 연관성 분석</h2>
          <p>연관업무 기록 조회 기능은 현재 개발 중입니다.</p>
        </WorkRecords>
      </DashboardContainer>
    </ContentArea>
  </RelatedWorkContainer>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  RelatedWorkContainer, ContentArea, DashboardContainer, WorkRecords
} from './RelatedWork.styled.js';

export default {
  name: 'RelatedWork',
  components: {
    SideBar,
    RelatedWorkContainer, ContentArea, DashboardContainer, WorkRecords
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