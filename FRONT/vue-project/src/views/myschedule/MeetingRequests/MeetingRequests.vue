<template>
  <Container>
    <SideBar />
    <ContentArea>
      <Title>협업/미팅 요청 관리</Title>
      <Tabs>
        <Tab :class="{ active: activeTab === 'received' }" @click="activeTab = 'received'">
          받은 요청 ({{ receivedRequests.length }})
        </Tab>
        <Tab :class="{ active: activeTab === 'sent' }" @click="activeTab = 'sent'">
          보낸 요청 ({{ sentRequests.length }})
        </Tab>
      </Tabs>
      <RequestList>
        <EmptyState v-if="activeTab === 'received' && receivedRequests.length === 0">받은 미팅 요청이 없습니다.</EmptyState>
        <EmptyState v-if="activeTab === 'sent' && sentRequests.length === 0">보낸 미팅 요청이 없습니다.</EmptyState>
        <template v-if="activeTab === 'received'">
          <RequestCard v-for="request in receivedRequests" :key="request.id" :class="[request.status]">
            <RequestHeader>
              <RequestTitle>{{ request.title }}</RequestTitle>
              <RequestStatus>{{ getStatusText(request.status) }}</RequestStatus>
            </RequestHeader>
            <RequestDetails>
              <div>요청자: {{ request.from }}</div>
              <div>{{ formatDateTime(request.dateTime) }}</div>
              <div>소요 시간: {{ request.duration }}분</div>
              <div>{{ request.description }}</div>
            </RequestDetails>
            <RequestActions v-if="request.status === 'pending'">
              <AcceptBtn @click="acceptRequest(request.id)">수락</AcceptBtn>
              <RejectBtn @click="rejectRequest(request.id)">거절</RejectBtn>
              <SuggestBtn @click="suggestAlternative(request.id)">대안 제시</SuggestBtn>
            </RequestActions>
          </RequestCard>
        </template>
        <template v-if="activeTab === 'sent'">
          <RequestCard v-for="request in sentRequests" :key="request.id" :class="[request.status]">
            <RequestHeader>
              <RequestTitle>{{ request.title }}</RequestTitle>
              <RequestStatus>{{ getStatusText(request.status) }}</RequestStatus>
            </RequestHeader>
            <RequestDetails>
              <div>수신자: {{ request.to }}</div>
              <div>{{ formatDateTime(request.dateTime) }}</div>
              <div>소요 시간: {{ request.duration }}분</div>
              <div>{{ request.description }}</div>
            </RequestDetails>
            <RequestActions v-if="request.status === 'pending'">
              <CancelBtn @click="cancelRequest(request.id)">취소</CancelBtn>
              <EditBtn @click="editRequest(request.id)">수정</EditBtn>
            </RequestActions>
            <div v-if="request.status === 'accepted'">
              <div>{{ request.to }}님이 미팅을 수락했습니다.</div>
            </div>
            <div v-if="request.status === 'rejected'">
              <div>{{ request.to }}님이 미팅을 거절했습니다.</div>
              <div v-if="request.rejectReason">사유: {{ request.rejectReason }}</div>
            </div>
          </RequestCard>
        </template>
      </RequestList>
      <ActionBar>
        <NewRequestBtn @click="showNewRequestForm">+ 새 미팅 요청</NewRequestBtn>
      </ActionBar>
      <ModalOverlay v-if="showRequestForm">
        <ModalContent>
          <ModalHeader>
            <h2>{{ isEditMode ? '미팅 요청 수정' : '새 미팅 요청' }}</h2>
            <CloseBtn @click="closeRequestForm">×</CloseBtn>
          </ModalHeader>
          <ModalBody>
            <form @submit.prevent="submitRequest">
              <FormGroup>
                <label>제목</label>
                <input type="text" v-model="requestForm.title" required />
              </FormGroup>
              <FormGroup>
                <label>수신자</label>
                <select v-model="requestForm.to" required>
                  <option value="">수신자 선택</option>
                  <option v-for="member in teamMembers" :key="member.id" :value="member.name">
                    {{ member.position }} {{ member.name }}
                  </option>
                </select>
              </FormGroup>
              <FormGroup>
                <label>날짜 및 시간</label>
                <input type="datetime-local" v-model="requestForm.dateTime" required />
              </FormGroup>
              <FormGroup>
                <label>소요 시간 (분)</label>
                <input type="number" v-model="requestForm.duration" min="15" step="15" required />
              </FormGroup>
              <FormGroup>
                <label>설명</label>
                <textarea v-model="requestForm.description" rows="3"></textarea>
              </FormGroup>
              <FormActions>
                <CancelBtn type="button" @click="closeRequestForm">취소</CancelBtn>
                <SubmitBtn type="submit">{{ isEditMode ? '수정' : '요청' }}</SubmitBtn>
              </FormActions>
            </form>
          </ModalBody>
        </ModalContent>
      </ModalOverlay>
      <ModalOverlay v-if="showAlternativeForm">
        <ModalContent>
          <ModalHeader>
            <h2>대안 시간 제시</h2>
            <CloseBtn @click="closeAlternativeForm">×</CloseBtn>
          </ModalHeader>
          <ModalBody>
            <p>{{ selectedRequest ? selectedRequest.from : '' }}님에게 다른 시간을 제안합니다.</p>
            <form @submit.prevent="submitAlternative">
              <FormGroup>
                <label>대안 날짜 및 시간</label>
                <input type="datetime-local" v-model="alternativeForm.dateTime" required />
              </FormGroup>
              <FormGroup>
                <label>메시지</label>
                <textarea v-model="alternativeForm.message" rows="3"></textarea>
              </FormGroup>
              <FormActions>
                <CancelBtn type="button" @click="closeAlternativeForm">취소</CancelBtn>
                <SubmitBtn type="submit">제안하기</SubmitBtn>
              </FormActions>
            </form>
          </ModalBody>
        </ModalContent>
      </ModalOverlay>
    </ContentArea>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, ContentArea, Title, Tabs, Tab, RequestList, RequestCard, RequestHeader, RequestTitle, RequestStatus, RequestDetails, RequestActions, AcceptBtn, RejectBtn, SuggestBtn, CancelBtn, EditBtn, NewRequestBtn, ActionBar, ModalOverlay, ModalContent, ModalHeader, ModalBody, FormGroup, FormActions, SubmitBtn, CloseBtn, EmptyState
} from './MeetingRequests.styled.js';

export default {
  name: 'MeetingRequests',
  components: {
    SideBar,
    Container, ContentArea, Title, Tabs, Tab, RequestList, RequestCard, RequestHeader, RequestTitle, RequestStatus, RequestDetails, RequestActions, AcceptBtn, RejectBtn, SuggestBtn, CancelBtn, EditBtn, NewRequestBtn, ActionBar, ModalOverlay, ModalContent, ModalHeader, ModalBody, FormGroup, FormActions, SubmitBtn, CloseBtn, EmptyState
  },
  data() {
    return {
      activeTab: 'received',
      showRequestForm: false,
      showAlternativeForm: false,
      isEditMode: false,
      selectedRequestId: null,
      requestForm: {
        title: '',
        to: '',
        dateTime: '',
        duration: 30,
        description: ''
      },
      alternativeForm: {
        dateTime: '',
        message: ''
      },
      teamMembers: [
        { id: 1, name: '김부장', position: '부장' },
        { id: 2, name: '이차장', position: '차장' },
        { id: 3, name: '박과장', position: '과장' },
        { id: 4, name: '최대리', position: '대리' },
        { id: 5, name: '정사원', position: '사원' }
      ],
      receivedRequests: [
        { id: 1, title: '프로젝트 킥오프 미팅', from: '김부장', dateTime: '2023-06-20T14:00', duration: 60, description: '신규 프로젝트 킥오프 미팅입니다. 참석 부탁드립니다.', status: 'pending' },
        { id: 2, title: '디자인 검토 회의', from: '최대리', dateTime: '2023-06-21T10:00', duration: 30, description: 'UI 디자인 검토 회의입니다.', status: 'accepted' },
        { id: 3, title: '코드 리뷰 세션', from: '박과장', dateTime: '2023-06-22T15:00', duration: 45, description: '백엔드 코드 리뷰 세션입니다.', status: 'rejected', rejectReason: '해당 시간에 다른 일정이 있습니다.' }
      ],
      sentRequests: [
        { id: 4, title: '기능 구현 논의', to: '이차장', dateTime: '2023-06-23T11:00', duration: 45, description: '신규 기능 구현 방향 논의', status: 'pending' },
        { id: 5, title: '버그 수정 회의', to: '정사원', dateTime: '2023-06-24T13:30', duration: 30, description: '발견된 버그 수정 방안 논의', status: 'accepted' }
      ]
    };
  },
  computed: {
    selectedRequest() {
      if (!this.selectedRequestId) return null;
      return this.receivedRequests.find(r => r.id === this.selectedRequestId);
    }
  },
  methods: {
    getStatusText(status) {
      switch (status) {
        case 'pending': return '대기중';
        case 'accepted': return '수락됨';
        case 'rejected': return '거절됨';
        default: return '';
      }
    },
    formatDateTime(dateTimeString) {
      const date = new Date(dateTimeString);
      return date.toLocaleString('ko-KR', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true
      });
    },
    showNewRequestForm() {
      this.isEditMode = false;
      this.resetRequestForm();
      this.showRequestForm = true;
    },
    closeRequestForm() {
      this.showRequestForm = false;
      this.resetRequestForm();
    },
    resetRequestForm() {
      this.requestForm = {
        title: '',
        to: '',
        dateTime: '',
        duration: 30,
        description: ''
      };
    },
    submitRequest() {
      if (this.isEditMode) {
        const index = this.sentRequests.findIndex(r => r.id === this.selectedRequestId);
        if (index !== -1) {
          this.sentRequests[index] = {
            ...this.sentRequests[index],
            title: this.requestForm.title,
            to: this.requestForm.to,
            dateTime: this.requestForm.dateTime,
            duration: this.requestForm.duration,
            description: this.requestForm.description
          };
        }
      } else {
        const newId = Math.max(...this.sentRequests.map(r => r.id), ...this.receivedRequests.map(r => r.id)) + 1;
        this.sentRequests.push({
          id: newId,
          title: this.requestForm.title,
          to: this.requestForm.to,
          dateTime: this.requestForm.dateTime,
          duration: this.requestForm.duration,
          description: this.requestForm.description,
          status: 'pending'
        });
      }
      this.closeRequestForm();
    },
    editRequest(requestId) {
      this.isEditMode = true;
      this.selectedRequestId = requestId;
      const request = this.sentRequests.find(r => r.id === requestId);
      if (request) {
        this.requestForm = {
          title: request.title,
          to: request.to,
          dateTime: request.dateTime,
          duration: request.duration,
          description: request.description
        };
        this.showRequestForm = true;
      }
    },
    cancelRequest(requestId) {
      if (confirm('이 미팅 요청을 취소하시겠습니까?')) {
        this.sentRequests = this.sentRequests.filter(r => r.id !== requestId);
      }
    },
    acceptRequest(requestId) {
      const index = this.receivedRequests.findIndex(r => r.id === requestId);
      if (index !== -1) {
        this.receivedRequests[index].status = 'accepted';
      }
    },
    rejectRequest(requestId) {
      const index = this.receivedRequests.findIndex(r => r.id === requestId);
      if (index !== -1) {
        const reason = prompt('거절 사유를 입력해주세요 (선택사항):');
        this.receivedRequests[index].status = 'rejected';
        if (reason) {
          this.receivedRequests[index].rejectReason = reason;
        }
      }
    },
    suggestAlternative(requestId) {
      this.selectedRequestId = requestId;
      this.alternativeForm = {
        dateTime: '',
        message: '해당 시간에는 일정이 있어 다른 시간을 제안드립니다.'
      };
      this.showAlternativeForm = true;
    },
    closeAlternativeForm() {
      this.showAlternativeForm = false;
      this.selectedRequestId = null;
    },
    submitAlternative() {
      alert('대안 시간이 제안되었습니다.');
      this.closeAlternativeForm();
      const index = this.receivedRequests.findIndex(r => r.id === this.selectedRequestId);
      if (index !== -1) {
        this.receivedRequests[index].status = 'pending';
        this.receivedRequests[index].alternativeProposed = true;
      }
    }
  }
};
</script>

<!-- 스타일은 vue-styled-components로 분리됨 -->
 