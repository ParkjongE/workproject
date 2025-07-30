<template>
  <Container>
    <SideBar />
    <MainContent>
      <TopActions>
        <ActionBtn>+ 구성원 추가</ActionBtn>
      </TopActions>
      <Title>구성원 관리</Title>
      <SearchSortRow>
        <SearchBox>
          <SearchInput type="text" placeholder="이름 또는 부서 검색" v-model="searchQuery" />
          <SearchBtn @click="searchMembers" class="search-btn">
            <font-awesome-icon :icon="['fas', 'search']" />
          </SearchBtn>
        </SearchBox>
        <SortGroupBox>
          <span @click="sortBy('department')">부서별</span>
          <span @click="sortBy('name')">이름순</span>
          <span @click="sortBy('position')">직급순</span>
        </SortGroupBox>
      </SearchSortRow>
      <MemberList>
        <MemberCard v-for="(member, index) in filteredMembers" :key="index">
          <MemberPhoto>
            <span>{{ member.photoText }}</span>
          </MemberPhoto>
          <h3>{{ member.name }}</h3>
          <MemberInfo>
            <MemberInfoRow>
              <Label>부서</Label>
              <span>{{ member.department }}</span>
            </MemberInfoRow>
            <MemberInfoRow>
              <Label>팀</Label>
              <span>{{ member.team }}</span>
            </MemberInfoRow>
            <MemberInfoRow>
              <Label>직급</Label>
              <span>{{ member.position }}</span>
            </MemberInfoRow>
            <MemberInfoRow>
              <Label>연락처</Label>
              <span>{{ member.contact }}</span>
            </MemberInfoRow>
            <MemberInfoRow>
              <Label>이메일</Label>
              <span>{{ member.email }}</span>
            </MemberInfoRow>
          </MemberInfo>
          <div style="display: flex; gap: 8px; margin-top: 12px;">
            <MemberStatus>{{ member.status }}</MemberStatus>
            <EditBtn @click="editMember(member)">수정</EditBtn>
          </div>
        </MemberCard>
      </MemberList>
    </MainContent>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, MainContent, TopActions, ActionBtn, Title, SearchSortRow, SearchBox, SearchInput, SearchBtn, SortGroupBox, MemberList, MemberCard, MemberPhoto, MemberInfo, MemberInfoRow, Label, MemberStatus, EditBtn
} from './OrgMember.styled.js';

export default {
  name: 'OrgMember',
  components: {
    SideBar,
    Container, MainContent, TopActions, ActionBtn, Title, SearchSortRow, SearchBox, SearchInput, SearchBtn, SortGroupBox, MemberList, MemberCard, MemberPhoto, MemberInfo, MemberInfoRow, Label, MemberStatus, EditBtn
  },
  data() {
    return {
      searchQuery: '',
      sortBy: 'name',
      members: [
        {
          id: 1,
          name: '김철수',
          position: '팀장',
          department: '개발팀',
          email: 'kim@company.com',
          phone: '010-1234-5678',
          status: '재직중'
        },
        {
          id: 2,
          name: '이영희',
          position: '사원',
          department: '디자인팀',
          email: 'lee@company.com',
          phone: '010-2345-6789',
          status: '재직중'
        },
        {
          id: 3,
          name: '박민수',
          position: '대리',
          department: '마케팅팀',
          email: 'park@company.com',
          phone: '010-3456-7890',
          status: '재직중'
        },
        {
          id: 4,
          name: '최지영',
          position: '사원',
          department: '개발팀',
          email: 'choi@company.com',
          phone: '010-4567-8901',
          status: '재직중'
        },
        {
          id: 5,
          name: '정현우',
          position: '팀장',
          department: '영업팀',
          email: 'jung@company.com',
          phone: '010-5678-9012',
          status: '재직중'
        },
        {
          id: 6,
          name: '한소영',
          position: '사원',
          department: '인사팀',
          email: 'han@company.com',
          phone: '010-6789-0123',
          status: '재직중'
        }
      ]
    };
  },
  computed: {
    filteredMembers() {
      return this.members.filter(member => 
        member.name.toLowerCase().includes(this.searchQuery.toLowerCase()) ||
        member.department.toLowerCase().includes(this.searchQuery.toLowerCase()) ||
        member.position.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
    sortedMembers() {
      return [...this.filteredMembers].sort((a, b) => {
        switch (this.sortBy) {
          case 'name':
            return a.name.localeCompare(b.name);
          case 'department':
            return a.department.localeCompare(b.department);
          case 'position':
            return a.position.localeCompare(b.position);
          default:
            return 0;
        }
      });
    }
  },
  methods: {
    editMember(memberId) {
      console.log('편집할 멤버 ID:', memberId);
      // 편집 로직 구현
    }
  }
};
</script>

<style scoped>
/* CSS는 별도 파일로 분리됨 */
</style> 