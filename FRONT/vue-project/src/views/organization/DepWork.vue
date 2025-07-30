<template>
  <Container>
    <SideBar />
    <MainContent>
      <h1>부서별 업무 현황</h1>
      <HeaderContainer>
        <DepartmentSelect v-model="selectedDepartment">
          <option value="all">전체 부서</option>
          <option value="sales">영업부</option>
          <option value="tech">기술부</option>
          <option value="marketing">마케팅부</option>
          <option value="hr">인사부</option>
          <option value="management">경영부</option>
        </DepartmentSelect>
      </HeaderContainer>
      <div style="overflow-x: auto;">
        <Table>
          <thead>
            <tr>
              <Th>부서명</Th>
              <Th>팀명</Th>
              <Th>팀장</Th>
              <Th>인원 수</Th>
              <Th>주요 담당 업무</Th>
              <Th>업무명</Th>
              <Th>업무 담당자</Th>
              <Th>진행 상태</Th>
              <Th>우선순위</Th>
              <Th>시작일</Th>
              <Th>마감일</Th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(task, index) in filteredTasks" :key="index">
              <Td>{{ task.department }}</Td>
              <Td>{{ task.team }}</Td>
              <Td>{{ task.teamLeader }}</Td>
              <Td>{{ task.memberCount }}</Td>
              <Td>{{ task.mainTask }}</Td>
              <Td>{{ task.taskName }}</Td>
              <Td>{{ task.assignee }}</Td>
              <Td>
                <BadgeContainer>
                  <StatusBadge :status="task.status">{{ task.status }}</StatusBadge>
                </BadgeContainer>
              </Td>
              <Td>
                <BadgeContainer>
                  <PriorityBadge :class="['priority-badge', getPriorityClass(task.priority)]">{{ task.priority }}</PriorityBadge>
                </BadgeContainer>
              </Td>
              <Td>{{ task.startDate }}</Td>
              <Td>{{ task.dueDate }}</Td>
            </tr>
          </tbody>
        </Table>
      </div>
    </MainContent>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, MainContent, HeaderContainer, DepartmentSelect, Table, Th, Td, BadgeContainer, StatusBadge, PriorityBadge
} from './DepWork.styled.js';

export default {
  name: 'DepartmentWork',
  components: {
    SideBar,
    Container, MainContent, HeaderContainer, DepartmentSelect, Table, Th, Td, BadgeContainer, StatusBadge, PriorityBadge
  },
  data() {
    return {
      selectedDepartment: 'all',
      tasks: [
        {
          department: '영업부',
          team: '영업1팀',
          teamLeader: '김영업',
          memberCount: 5,
          mainTask: '국내 영업',
          taskName: '신규 거래처 발굴',
          assignee: '이영업',
          status: '진행중',
          priority: '높음',
          startDate: '2024-07-01',
          dueDate: '2024-07-31'
        },
        {
          department: '영업부',
          team: '영업2팀',
          teamLeader: '박영업',
          memberCount: 4,
          mainTask: '해외 영업',
          taskName: '해외 전시회 참가',
          assignee: '최영업',
          status: '완료',
          priority: '중간',
          startDate: '2024-06-15',
          dueDate: '2024-07-15'
        },
        {
          department: '기술부',
          team: '기술1팀',
          teamLeader: '김기술',
          memberCount: 6,
          mainTask: '소프트웨어 개발',
          taskName: '신규 기능 개발',
          assignee: '이기술',
          status: '지연',
          priority: '긴급',
          startDate: '2024-06-01',
          dueDate: '2024-07-15'
        },
        {
          department: '기술부',
          team: '기술2팀',
          teamLeader: '박기술',
          memberCount: 5,
          mainTask: '하드웨어 개발',
          taskName: '신제품 테스트',
          assignee: '최기술',
          status: '진행중',
          priority: '중간',
          startDate: '2024-07-10',
          dueDate: '2024-08-10'
        },
        {
          department: '마케팅부',
          team: '마케팅팀',
          teamLeader: '김마케팅',
          memberCount: 4,
          mainTask: '마케팅 전략 수립',
          taskName: '신제품 홍보 전략',
          assignee: '이마케팅',
          status: '진행중',
          priority: '높음',
          startDate: '2024-07-05',
          dueDate: '2024-07-25'
        },
        {
          department: '인사부',
          team: '인사팀',
          teamLeader: '김인사',
          memberCount: 3,
          mainTask: '인사 관리',
          taskName: '신입 사원 채용',
          assignee: '이인사',
          status: '완료',
          priority: '중간',
          startDate: '2024-06-20',
          dueDate: '2024-07-10'
        },
        {
          department: '경영부',
          team: '경영지원팀',
          teamLeader: '김경영',
          memberCount: 4,
          mainTask: '경영 지원',
          taskName: '분기 보고서 작성',
          assignee: '이경영',
          status: '진행중',
          priority: '높음',
          startDate: '2024-07-01',
          dueDate: '2024-07-20'
        }
      ]
    };
  },
  computed: {
    filteredTasks() {
      if (this.selectedDepartment === 'all') {
        return this.tasks;
      }
      
      const departmentMap = {
        'sales': '영업부',
        'tech': '기술부',
        'marketing': '마케팅부',
        'hr': '인사부',
        'management': '경영부'
      };
      
      return this.tasks.filter(task => task.department === departmentMap[this.selectedDepartment]);
    }
  },
  methods: {
    getStatusClass(status) {
      switch (status) {
        case '진행중': return 'status-progress';
        case '완료': return 'status-completed';
        case '지연': return 'status-delayed';
        default: return '';
      }
    },
    getPriorityClass(priority) {
      switch (priority) {
        case '긴급': return 'priority-urgent';
        case '높음': return 'priority-urgent';
        case '중간': return 'priority-medium';
        case '낮음': return 'priority-low';
        default: return '';
      }
    }
  }
};
</script>

<style scoped>
/* CSS는 별도 파일로 분리됨 */
</style> 