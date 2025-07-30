<template>
  <WorkRecordsContainer>
    <SideBar />
    <ContentArea>
      <h1>내 업무 기록</h1>
      
      <FilterSection>
        <DateFilter>
          <label>기간:</label>
          <select v-model="dateFilter">
            <option value="today">오늘</option>
            <option value="week">이번 주</option>
            <option value="month">이번 달</option>
            <option value="custom">직접 설정</option>
          </select>
          
          <div v-if="dateFilter === 'custom'" class="custom-date-range">
            <input type="date" v-model="startDate">
            <span>~</span>
            <input type="date" v-model="endDate">
            <button class="apply-btn" @click="applyDateFilter">적용</button>
          </div>
        </DateFilter>
        
        <CategoryFilter>
          <label>카테고리:</label>
          <select v-model="categoryFilter">
            <option value="all">전체</option>
            <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
          </select>
        </CategoryFilter>
        
        <SearchBox>
          <input type="text" v-model="searchKeyword" placeholder="검색어 입력...">
          <button class="search-btn" @click="searchRecords">검색</button>
        </SearchBox>
      </FilterSection>
      
      <StatsSection>
        <StatCard>
          <StatTitle>총 업무 시간</StatTitle>
          <StatValue>{{ getTotalWorkHours() }}시간</StatValue>
        </StatCard>
        <StatCard>
          <StatTitle>완료된 업무</StatTitle>
          <StatValue>{{ getCompletedTasksCount() }}개</StatValue>
        </StatCard>
        <StatCard>
          <StatTitle>진행중인 업무</StatTitle>
          <StatValue>{{ getInProgressTasksCount() }}개</StatValue>
        </StatCard>
      </StatsSection>
      
      <RecordsSection>
        <SectionHeader>
          <h2>업무 기록</h2>
          <AddRecordBtn @click="showAddRecordForm">+ 새 업무 기록</AddRecordBtn>
        </SectionHeader>
        
        <EmptyState v-if="filteredRecords.length === 0">
          업무 기록이 없습니다.
        </EmptyState>
        
        <RecordsList v-else>
          <RecordCard 
            v-for="record in filteredRecords" 
            :key="record.id"
            :class="{ 'completed': record.status === 'completed' }"
          >
            <RecordHeader>
              <RecordTitle>{{ record.title }}</RecordTitle>
              <RecordStatus :class="record.status">
                {{ getStatusText(record.status) }}
              </RecordStatus>
            </RecordHeader>
            <RecordDetails>
              <RecordDate>{{ formatDate(record.date) }}</RecordDate>
              <RecordTime>{{ record.startTime }} ~ {{ record.endTime }}</RecordTime>
              <RecordDuration>{{ calculateDuration(record.startTime, record.endTime) }}시간</RecordDuration>
              <RecordCategory>{{ record.category }}</RecordCategory>
            </RecordDetails>
            <RecordDescription>
              {{ record.description }}
            </RecordDescription>
            <RecordActions>
              <EditBtn @click="editRecord(record.id)">수정</EditBtn>
              <DeleteBtn @click="deleteRecord(record.id)">삭제</DeleteBtn>
              <CompleteBtn 
                v-if="record.status === 'in-progress'" 
                @click="completeRecord(record.id)"
              >
                완료
              </CompleteBtn>
            </RecordActions>
          </RecordCard>
        </RecordsList>
      </RecordsSection>
      
      <!-- 업무 기록 추가/수정 모달 -->
      <ModalOverlay v-if="showRecordForm">
        <ModalContent>
          <ModalHeader>
            <h2>{{ isEditMode ? '업무 기록 수정' : '새 업무 기록' }}</h2>
            <CloseBtn @click="closeRecordForm">×</CloseBtn>
          </ModalHeader>
          <ModalBody>
            <form @submit.prevent="saveRecord">
              <FormGroup>
                <label>제목</label>
                <input type="text" v-model="recordForm.title" required>
              </FormGroup>
              <FormGroup>
                <label>날짜</label>
                <input type="date" v-model="recordForm.date" required>
              </FormGroup>
              <FormRow>
                <FormGroup half>
                  <label>시작 시간</label>
                  <input type="time" v-model="recordForm.startTime" required>
                </FormGroup>
                <FormGroup half>
                  <label>종료 시간</label>
                  <input type="time" v-model="recordForm.endTime" required>
                </FormGroup>
              </FormRow>
              <FormGroup>
                <label>카테고리</label>
                <select v-model="recordForm.category" required>
                  <option value="">카테고리 선택</option>
                  <option v-for="category in categories" :key="category" :value="category">
                    {{ category }}
                  </option>
                </select>
              </FormGroup>
              <FormGroup>
                <label>상태</label>
                <select v-model="recordForm.status" required>
                  <option value="in-progress">진행중</option>
                  <option value="completed">완료</option>
                </select>
              </FormGroup>
              <FormGroup>
                <label>설명</label>
                <textarea v-model="recordForm.description" rows="3"></textarea>
              </FormGroup>
              <FormActions>
                <CancelBtn @click="closeRecordForm">취소</CancelBtn>
                <SaveBtn>저장</SaveBtn>
              </FormActions>
            </form>
          </ModalBody>
        </ModalContent>
      </ModalOverlay>
    </ContentArea>
  </WorkRecordsContainer>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  WorkRecordsContainer, ContentArea, FilterSection, DateFilter, CategoryFilter, SearchBox, StatsSection, StatCard, StatTitle, StatValue, RecordsSection, SectionHeader, AddRecordBtn, EmptyState, RecordsList, RecordCard, RecordHeader, RecordTitle, RecordStatus, RecordDetails, RecordDate, RecordTime, RecordDuration, RecordCategory, RecordDescription, RecordActions, EditBtn, DeleteBtn, CompleteBtn, ModalOverlay, ModalContent, ModalHeader, CloseBtn, ModalBody, FormGroup, FormRow, FormActions, CancelBtn, SaveBtn
} from './WorkRecords.styled.js';

export default {
  name: 'WorkRecords',
  components: {
    SideBar,
    WorkRecordsContainer, ContentArea, FilterSection, DateFilter, CategoryFilter, SearchBox, StatsSection, StatCard, StatTitle, StatValue, RecordsSection, SectionHeader, AddRecordBtn, EmptyState, RecordsList, RecordCard, RecordHeader, RecordTitle, RecordStatus, RecordDetails, RecordDate, RecordTime, RecordDuration, RecordCategory, RecordDescription, RecordActions, EditBtn, DeleteBtn, CompleteBtn, ModalOverlay, ModalContent, ModalHeader, CloseBtn, ModalBody, FormGroup, FormRow, FormActions, CancelBtn, SaveBtn
  },
  data() {
    return {
      dateFilter: 'week',
      categoryFilter: 'all',
      searchKeyword: '',
      startDate: this.getFormattedDate(new Date(new Date().setDate(new Date().getDate() - 7))),
      endDate: this.getFormattedDate(new Date()),
      showRecordForm: false,
      isEditMode: false,
      selectedRecordId: null,
      recordForm: {
        title: '',
        date: '',
        startTime: '',
        endTime: '',
        category: '',
        status: 'in-progress',
        description: ''
      },
      categories: [
        '개발',
        '회의',
        '문서작업',
        '기획',
        '디자인',
        '테스트',
        '기타'
      ],
      // 예시 데이터
      workRecords: [
        {
          id: 1,
          title: '로그인 기능 개발',
          date: '2023-06-15',
          startTime: '09:00',
          endTime: '12:30',
          category: '개발',
          status: 'completed',
          description: '로그인 및 인증 기능 개발 완료'
        },
        {
          id: 2,
          title: '팀 주간 회의',
          date: '2023-06-15',
          startTime: '14:00',
          endTime: '15:30',
          category: '회의',
          status: 'completed',
          description: '이번 주 진행 상황 및 다음 주 계획 논의'
        },
        {
          id: 3,
          title: '대시보드 UI 개발',
          date: '2023-06-16',
          startTime: '09:00',
          endTime: '18:00',
          category: '개발',
          status: 'in-progress',
          description: '사용자 대시보드 UI 컴포넌트 개발 중'
        },
        {
          id: 4,
          title: '요구사항 명세서 작성',
          date: '2023-06-17',
          startTime: '10:00',
          endTime: '12:00',
          category: '문서작업',
          status: 'in-progress',
          description: '신규 기능 요구사항 명세서 작성 중'
        },
        {
          id: 5,
          title: '버그 수정',
          date: '2023-06-18',
          startTime: '13:00',
          endTime: '15:30',
          category: '개발',
          status: 'completed',
          description: '사용자 프로필 페이지 버그 수정'
        }
      ]
    };
  },
  computed: {
    filteredRecords() {
      let records = [...this.workRecords];
      
      // 날짜 필터링
      if (this.dateFilter === 'today') {
        const today = this.getFormattedDate(new Date());
        records = records.filter(record => record.date === today);
      } else if (this.dateFilter === 'week') {
        const oneWeekAgo = new Date();
        oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
        records = records.filter(record => new Date(record.date) >= oneWeekAgo);
      } else if (this.dateFilter === 'month') {
        const oneMonthAgo = new Date();
        oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);
        records = records.filter(record => new Date(record.date) >= oneMonthAgo);
      } else if (this.dateFilter === 'custom') {
        records = records.filter(record => 
          new Date(record.date) >= new Date(this.startDate) && 
          new Date(record.date) <= new Date(this.endDate)
        );
      }
      
      // 카테고리 필터링
      if (this.categoryFilter !== 'all') {
        records = records.filter(record => record.category === this.categoryFilter);
      }
      
      // 검색어 필터링
      if (this.searchKeyword) {
        const keyword = this.searchKeyword.toLowerCase();
        records = records.filter(record => 
          record.title.toLowerCase().includes(keyword) ||
          record.description.toLowerCase().includes(keyword)
        );
      }
      
      // 날짜 기준 내림차순 정렬
      return records.sort((a, b) => new Date(b.date) - new Date(a.date));
    }
  },
  methods: {
    getFormattedDate(date) {
      return date.toISOString().split('T')[0];
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ko-KR', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' });
    },
    calculateDuration(startTime, endTime) {
      const start = new Date(`2000-01-01T${startTime}`);
      const end = new Date(`2000-01-01T${endTime}`);
      const diffMs = end - start;
      const diffHours = diffMs / (1000 * 60 * 60);
      return diffHours.toFixed(1);
    },
    getStatusText(status) {
      return status === 'completed' ? '완료' : '진행중';
    },
    getTotalWorkHours() {
      let total = 0;
      this.filteredRecords.forEach(record => {
        total += parseFloat(this.calculateDuration(record.startTime, record.endTime));
      });
      return total.toFixed(1);
    },
    getCompletedTasksCount() {
      return this.filteredRecords.filter(record => record.status === 'completed').length;
    },
    getInProgressTasksCount() {
      return this.filteredRecords.filter(record => record.status === 'in-progress').length;
    },
    applyDateFilter() {
      // 날짜 필터 적용 로직
      console.log(`날짜 필터 적용: ${this.startDate} ~ ${this.endDate}`);
    },
    searchRecords() {
      // 검색 로직
      console.log(`검색어: ${this.searchKeyword}`);
    },
    showAddRecordForm() {
      this.isEditMode = false;
      this.resetRecordForm();
      this.recordForm.date = this.getFormattedDate(new Date());
      this.showRecordForm = true;
    },
    closeRecordForm() {
      this.showRecordForm = false;
      this.resetRecordForm();
    },
    resetRecordForm() {
      this.recordForm = {
        title: '',
        date: '',
        startTime: '',
        endTime: '',
        category: '',
        status: 'in-progress',
        description: ''
      };
    },
    saveRecord() {
      if (this.isEditMode) {
        // 기존 기록 수정
        const index = this.workRecords.findIndex(r => r.id === this.selectedRecordId);
        if (index !== -1) {
          this.workRecords[index] = {
            ...this.workRecords[index],
            title: this.recordForm.title,
            date: this.recordForm.date,
            startTime: this.recordForm.startTime,
            endTime: this.recordForm.endTime,
            category: this.recordForm.category,
            status: this.recordForm.status,
            description: this.recordForm.description
          };
        }
      } else {
        // 새 기록 추가
        const newId = Math.max(...this.workRecords.map(r => r.id)) + 1;
        this.workRecords.push({
          id: newId,
          title: this.recordForm.title,
          date: this.recordForm.date,
          startTime: this.recordForm.startTime,
          endTime: this.recordForm.endTime,
          category: this.recordForm.category,
          status: this.recordForm.status,
          description: this.recordForm.description
        });
      }
      
      this.closeRecordForm();
    },
    editRecord(recordId) {
      this.isEditMode = true;
      this.selectedRecordId = recordId;
      
      const record = this.workRecords.find(r => r.id === recordId);
      if (record) {
        this.recordForm = {
          title: record.title,
          date: record.date,
          startTime: record.startTime,
          endTime: record.endTime,
          category: record.category,
          status: record.status,
          description: record.description
        };
        
        this.showRecordForm = true;
      }
    },
    deleteRecord(recordId) {
      if (confirm('이 업무 기록을 삭제하시겠습니까?')) {
        // 실제 구현에서는 API 호출
        this.workRecords = this.workRecords.filter(r => r.id !== recordId);
      }
    },
    completeRecord(recordId) {
      const index = this.workRecords.findIndex(r => r.id === recordId);
      if (index !== -1) {
        this.workRecords[index].status = 'completed';
      }
    }
  }
};
</script>
 