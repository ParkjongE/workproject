// 라이브러리 없이 전역 상태 관리 구현 
import { reactive, computed } from 'vue'

// 전역 상태 객체 생성
const scheduleState = reactive({
  // State
  selectedUser: null,
  teamMembers: [],
  isUserScheduleModalOpen: false,
  modalType: 'view', // 'view', 'edit', 'create'
  currentSchedule: null,
  isLoading: false,
  error: null
})

// Getters (computed로 구현)
const getSelectedUser = computed(() => scheduleState.selectedUser)
const getModalStatus = computed(() => scheduleState.isUserScheduleModalOpen)
const getModalType = computed(() => scheduleState.modalType)
const getTeamMembers = computed(() => scheduleState.teamMembers)

// Actions
const openUserScheduleModal = (user) => {
  scheduleState.selectedUser = user
  scheduleState.isUserScheduleModalOpen = true
  scheduleState.modalType = 'view'
  console.log('Store: 사용자 스케줄 모달 열기', user)
}

const closeUserScheduleModal = () => {
  scheduleState.isUserScheduleModalOpen = false
  scheduleState.selectedUser = null
  scheduleState.currentSchedule = null
  console.log('Store: 모달 닫기')
}

const setTeamMembers = (members) => {
  scheduleState.teamMembers = members
  console.log('Store: 팀원 데이터 설정', members.length)
}

// Vue 3 Composition API를 활용한 전역 상태 관리
export const useScheduleStore = () => {
  return {
    // State
    selectedUser: computed(() => scheduleState.selectedUser),
    teamMembers: computed(() => scheduleState.teamMembers),
    isUserScheduleModalOpen: computed(() => scheduleState.isUserScheduleModalOpen),
    modalType: computed(() => scheduleState.modalType),
    currentSchedule: computed(() => scheduleState.currentSchedule),
    isLoading: computed(() => scheduleState.isLoading),
    error: computed(() => scheduleState.error),
    
    // Getters
    getSelectedUser,
    getModalStatus,
    getModalType,
    getTeamMembers,
    
    // Actions
    openUserScheduleModal,
    closeUserScheduleModal,
    setTeamMembers
  }
}
