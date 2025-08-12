<!-- Commute.vue -->
<template>
  <CommuteRoot>
    <SideBar />
    <MainContent>
      <Title>출퇴근 기록</Title>
      <Card>
        <CardTitle>오늘의 출퇴근</CardTitle>
        <TimeDisplay>
          <TimeItem>
            <TimeLabel>출근 시간</TimeLabel>
            <TimeValue>{{ checkInTime || '미출근' }}</TimeValue>
          </TimeItem>
          <TimeItem>
            <TimeLabel>퇴근 시간</TimeLabel>
            <TimeValue>{{ checkOutTime || '미퇴근' }}</TimeValue>
          </TimeItem>
        </TimeDisplay>
        <ButtonGroup>
          <CheckInBtn
            :class="{ loading: isLoading }"
            @click="checkIn"
            :disabled="isCheckedIn"
          >
            {{ isCheckedIn ? '출근 완료' : '출근하기' }}
          </CheckInBtn>
          <CheckOutBtn
            :class="{ loading: isLoading }"
            @click="checkOut"
            :disabled="!isCheckedIn || isCheckedOut"
          >
            {{ isCheckedOut ? '퇴근 완료' : '퇴근하기' }}
          </CheckOutBtn>
        </ButtonGroup>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>이번 주 근무 현황</CardTitle>
          <CorrectionRequestBtn @click="showCorrectionModal = true">
            <i class="fas fa-edit"></i>
            근태 정정요청
          </CorrectionRequestBtn>
        </CardHeader>
        <WeekGrid>
          <WeekDay
            v-for="(day, index) in weekData"
            :key="index"
            :class="{ today: day.isToday, weekend: day.isWeekend }"
          >
            <DayName>{{ day.name }}</DayName>
            <DayDate>{{ day.date }}</DayDate>
            <DayStatus :class="day.status">{{ day.statusText }}</DayStatus>
          </WeekDay>
        </WeekGrid>
      </Card>
      <Card>
        <CardTitle>근태 이력</CardTitle>
        <HistoryTable>
          <Table>
            <thead>
              <Tr>
                <Th>날짜</Th>
                <Th>출근 시간</Th>
                <Th>퇴근 시간</Th>
                <Th>근무 시간</Th>
                <Th>상태</Th>
              </Tr>
            </thead>
            <tbody>
              <Tr v-for="record in attendanceHistory" :key="record.date">
                <Td>{{ record.date }}</Td>
                <Td>{{ record.checkIn || '-' }}</Td>
                <Td>{{ record.checkOut || '-' }}</Td>
                <Td>{{ record.workHours || '-' }}</Td>
                <Td>
                  <StatusBadge :class="record.status">{{ record.statusText }}</StatusBadge>
                </Td>
              </Tr>
            </tbody>
          </Table>
        </HistoryTable>
      </Card>
    </MainContent>
    <!-- 근태 정정요청 모달 -->
    <ModalOverlay v-if="showCorrectionModal" @click="closeModal">
      <ModalContent @click.stop>
        <ModalHeader>
          <h3>근태 정정요청</h3>
          <ModalClose @click="closeModal">
            <i class="fas fa-times"></i>
          </ModalClose>
        </ModalHeader>
        <ModalBody>
          <FormGroup>
            <label>정정 날짜</label>
            <FormInput type="date" v-model="correctionDate" />
          </FormGroup>
          <FormGroup>
            <label>정정 유형</label>
            <RadioGroup>
              <RadioLabel>
                <RadioInput
                  type="radio"
                  value="checkin"
                  v-model="correctionType"
                />
                출근 시간만
              </RadioLabel>
              <RadioLabel>
                <RadioInput
                  type="radio"
                  value="checkout"
                  v-model="correctionType"
                />
                퇴근 시간만
              </RadioLabel>
              <RadioLabel>
                <RadioInput
                  type="radio"
                  value="both"
                  v-model="correctionType"
                />
                출퇴근 시간 모두
              </RadioLabel>
            </RadioGroup>
          </FormGroup>
          <FormGroup v-if="correctionType === 'checkin' || correctionType === 'both'">
            <label>출근 시간</label>
            <FormInput type="time" v-model="correctionCheckIn" />
          </FormGroup>
          <FormGroup v-if="correctionType === 'checkout' || correctionType === 'both'">
            <label>퇴근 시간</label>
            <FormInput type="time" v-model="correctionCheckOut" />
          </FormGroup>
          <FormGroup>
            <label>정정 사유</label>
            <FormTextarea
              v-model="correctionReason"
              placeholder="정정 사유를 입력해주세요"
              rows="3"
            />
          </FormGroup>
        </ModalBody>
        <ModalFooter>
          <CancelBtn @click="closeModal">취소</CancelBtn>
          <SubmitBtn
            @click="submitCorrection"
            :disabled="isSubmitting"
            :class="{ loading: isSubmitting }"
          >
            {{ isSubmitting ? '제출 중...' : '제출' }}
          </SubmitBtn>
        </ModalFooter>
      </ModalContent>
    </ModalOverlay>
  </CommuteRoot>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import * as styled from './Commute.styled.js'
import { useAttendanceStore } from '@/stores/attendance'

export default {
  name: 'CommuteView',
  components: {
    SideBar,
    ...styled
  },
  setup() {
    const attendanceStore = useAttendanceStore()
    
    return {
      attendanceStore
    }
  },
  beforeMount() {
    // 컴포넌트 마운트 전에 스토어 초기화 확인
    if (!this.attendanceStore) {
      console.error('Attendance store가 초기화되지 않았습니다.')
    }
  },
  data() {
    return {
      isLoading: false,
      showCorrectionModal: false,
      isSubmitting: false,
      correctionDate: '',
      correctionType: 'checkin',
      correctionCheckIn: '',
      correctionCheckOut: '',
      correctionReason: ''
    }
  },
  computed: {
    // 스토어에서 가져온 데이터
    checkInTime() {
      const todayAttendance = this.attendanceStore?.todayAttendance?.value
      console.log('[checkInTime computed] todayAttendance 전체:', todayAttendance)
      const time = todayAttendance?.checkInTime
      console.log('[checkInTime computed] time:', time, 'type:', typeof time)
      return (time && time !== 'null' && time !== '') ? time : null
    },
    checkOutTime() {
      const time = this.attendanceStore?.todayAttendance?.value?.checkOutTime
      console.log('[checkOutTime computed] time:', time, 'type:', typeof time)
      return (time && time !== 'null' && time !== '') ? time : null
    },
    isCheckedIn() {
      const checkInTime = this.checkInTime
      console.log('[isCheckedIn] checkInTime:', checkInTime, 'type:', typeof checkInTime, 'result:', !!checkInTime)
      return !!checkInTime
    },
    isCheckedOut() {
      const checkOutTime = this.checkOutTime
      console.log('[isCheckedOut] checkOutTime:', checkOutTime, 'type:', typeof checkOutTime, 'result:', !!checkOutTime)
      return !!checkOutTime
    },
    attendanceHistory() {
      return (this.attendanceStore?.attendanceRecords?.value || []).map(record => ({
        date: record.date,
        checkIn: record.checkInTime,
        checkOut: record.checkOutTime,
        workHours: record.workHours ? `${Math.floor(record.workHours)}시간 ${Math.round((record.workHours % 1) * 60)}분` : null,
        status: this.getStatusClass(record),
        statusText: this.getStatusText(record)
      }))
    },
    weekData() {
      const today = new Date()
      const dayOfWeek = today.getDay()
      const startOfWeek = new Date(today)
      startOfWeek.setDate(today.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1))
      
      const weekDays = ['일', '월', '화', '수', '목', '금', '토']
      const weekData = []
      
      // 🔴 디버깅 로그 추가
      console.log('[weekData] attendanceRecords:', this.attendanceStore?.attendanceRecords?.value)
      
      for (let i = 0; i < 7; i++) {
        const date = new Date(startOfWeek)
        date.setDate(startOfWeek.getDate() + i)
        
        const dateStr = date.toISOString().split('T')[0]
        const dayName = weekDays[date.getDay()]
        const dayDate = date.getDate()
        const isToday = dateStr === today.toISOString().split('T')[0]
        const isWeekend = date.getDay() === 0 || date.getDay() === 6
        
        // 해당 날짜의 출퇴근 기록 찾기 (로깅 추가)
        const record = (this.attendanceStore?.attendanceRecords?.value || []).find(r => r.date === dateStr)
        console.log(`[weekData] ${dateStr} 날짜 찾기:`, record)
        
        let status = 'no-record'
        let statusText = '기록없음'
        
        if (isWeekend) {
          status = 'weekend'
          statusText = '주말'
        } else if (record) {
          console.log(`[weekData] ${dateStr} 레코드 발견:`, record)
          if (record.checkInTime && record.checkOutTime) {
            status = 'checked-out'
            statusText = '퇴근'
          } else if (record.checkInTime) {
            status = 'checked-in'
            statusText = '출근'
          }
        }
        
        weekData.push({
          name: dayName,
          date: dayDate,
          status,
          statusText,
          isToday,
          isWeekend
        })
      }
      
      return weekData
    }
  },
  methods: {
    async checkIn() {
      this.isLoading = true
      try {
        await this.attendanceStore.checkIn()
        // 출퇴근 기록 새로고침
        await this.attendanceStore.loadTodayAttendance()
        await this.attendanceStore.loadAttendanceRecords()
      } catch (error) {
        console.error('출근 처리 실패:', error)
      } finally {
        this.isLoading = false
      }
    },
    async checkOut() {
      this.isLoading = true
      try {
        await this.attendanceStore.checkOut()
        // 출퇴근 기록 새로고침
        await this.attendanceStore.loadTodayAttendance()
        await this.attendanceStore.loadAttendanceRecords()
      } catch (error) {
        console.error('퇴근 처리 실패:', error)
      } finally {
        this.isLoading = false
      }
    },
    getStatusClass(record) {
      if (!record.checkInTime || !record.checkOutTime) return 'incomplete'
      
      const workHours = record.workHours
      if (workHours >= 9) return 'overtime'
      if (workHours >= 8) return 'normal'
      return 'late'
    },
    getStatusText(record) {
      if (!record.checkInTime || !record.checkOutTime) return '미완료'
      
      const workHours = record.workHours
      if (workHours >= 9) return '초과근무'
      if (workHours >= 8) return '정상'
      return '지각'
    },
    closeModal() {
      this.showCorrectionModal = false
      this.resetCorrectionForm()
    },
    resetCorrectionForm() {
      this.correctionDate = ''
      this.correctionType = 'checkin'
      this.correctionCheckIn = ''
      this.correctionCheckOut = ''
      this.correctionReason = ''
    },
    async submitCorrection() {
      if (!this.correctionDate || !this.correctionReason) {
        alert('날짜와 사유를 모두 입력해주세요.')
        return
      }
      if (this.correctionType === 'checkin' && !this.correctionCheckIn) {
        alert('출근 시간을 입력해주세요.')
        return
      }
      if (this.correctionType === 'checkout' && !this.correctionCheckOut) {
        alert('퇴근 시간을 입력해주세요.')
        return
      }
      if (this.correctionType === 'both' && (!this.correctionCheckIn || !this.correctionCheckOut)) {
        alert('출퇴근 시간을 모두 입력해주세요.')
        return
      }
      this.isSubmitting = true
      try {
        await new Promise(resolve => setTimeout(resolve, 1500))
        alert('근태 정정요청이 제출되었습니다.')
        this.closeModal()
      } catch (error) {
        console.error('정정요청 제출 실패:', error)
        alert('정정요청 제출에 실패했습니다.')
      } finally {
        this.isSubmitting = false
      }
    }
  },
  async mounted() {
    // 컴포넌트 마운트 시 출퇴근 데이터 초기화
    console.log('[mounted] 출퇴근 데이터 초기화 시작')
    await this.attendanceStore.initializeAttendance()
    console.log('[mounted] 출퇴근 데이터 초기화 완료')
    console.log('[mounted] todayAttendance:', this.attendanceStore.todayAttendance?.value)
    console.log('[mounted] isCheckedIn:', this.isCheckedIn)
    console.log('[mounted] isCheckedOut:', this.isCheckedOut)
  }
}
</script>

<!-- 스타일은 vue-styled-components로 분리됨 --> 