// 라이브러리 없이 전역 상태 관리 구현 (Pinia 대체)
import { reactive, computed } from 'vue'
import api from '@/plugins/axios'

// 전역 상태 객체 생성
const attendanceState = reactive({
  // State
  currentTime: '',
  todayAttendance: null,
  attendanceRecords: [],
  workType: 'office',
  weeklyWorkMinutes: 0,
  weeklyWorkGoal: 40 * 60, // 40시간을 분으로
  timerInterval: null
})

// Getters (computed로 구현)
const workingHours = computed(() => {
  if (!attendanceState.todayAttendance || !attendanceState.todayAttendance.checkInTime || !attendanceState.todayAttendance.checkOutTime) {
    return '-'
  }
  
  // 원본 ISO 시간 문자열이 있으면 그것을 사용, 없으면 포맷된 시간을 파싱
  let checkIn, checkOut
  
  if (attendanceState.todayAttendance.originalCheckInTime) {
    checkIn = new Date(attendanceState.todayAttendance.originalCheckInTime)
    checkOut = new Date(attendanceState.todayAttendance.originalCheckOutTime)
  } else {
    // 포맷된 시간 문자열을 파싱 (HH:MM 형식)
    const [checkInHours, checkInMinutes] = attendanceState.todayAttendance.checkInTime.split(':').map(Number)
    const [checkOutHours, checkOutMinutes] = attendanceState.todayAttendance.checkOutTime.split(':').map(Number)
    
    const today = new Date()
    checkIn = new Date(today.getFullYear(), today.getMonth(), today.getDate(), checkInHours, checkInMinutes)
    checkOut = new Date(today.getFullYear(), today.getMonth(), today.getDate(), checkOutHours, checkOutMinutes)
  }
  
  const diffMs = checkOut - checkIn
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))
  
  return `${diffHours}시간 ${diffMinutes}분`
})

const weeklyWorkPercentage = computed(() => {
  return Math.min((attendanceState.weeklyWorkMinutes / attendanceState.weeklyWorkGoal) * 100, 100)
})

const getWorkStatus = computed(() => {
  if (weeklyWorkPercentage.value >= 100) return '완료'
  if (weeklyWorkPercentage.value >= 80) return '진행중'
  return '시작'
})

const getWorkStatusClass = computed(() => {
  if (weeklyWorkPercentage.value >= 100) return 'completed'
  if (weeklyWorkPercentage.value >= 80) return 'in-progress'
  return 'not-started'
})

// Actions
const updateCurrentTime = () => {
  const now = new Date()
  attendanceState.currentTime = now.toLocaleTimeString('ko-KR', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

const checkIn = async () => {
  try {
    const userStr = localStorage.getItem('user')
    if (!userStr) {
      alert('로그인 정보가 없습니다. 다시 로그인 해주세요.')
      return
    }
    
    const user = JSON.parse(userStr)
    const user_id = user.id

    // 현재 시간을 ISO 문자열로 전송
    const now = new Date()
    const response = await api.post('/api/attendance/checkin', {
      user_id: parseInt(user_id),
      date: now.toISOString().split('T')[0],
      check_in_time: now.toISOString()
    })

    if (response.data.success) {
      // 출퇴근 기록 다시 로드하여 최신 데이터 반영
      await loadTodayAttendance()
      await loadAttendanceRecords()
      
      alert(`출근이 기록되었습니다: ${response.data.check_in_time}`)
      
      // 지각 체크
      const checkInTime = parseTimeToMinutes(attendanceState.todayAttendance.checkInTime)
      const standardStartTime = 9 * 60
      
      if (checkInTime > standardStartTime) {
        alert(`지각입니다. 표준 출근 시간은 09:00입니다. (${Math.floor((checkInTime - standardStartTime) / 60)}시간 ${(checkInTime - standardStartTime) % 60}분 지각)`)
      }
    } else {
      alert('출근 기록에 실패했습니다.')
    }
  } catch (error) {
    console.error('[checkIn] 오류:', error)
    if (error.response?.status === 409) {
      alert('이미 출근 기록이 있습니다.')
    } else {
      alert('출근 기록 중 오류가 발생했습니다.')
    }
  }
}

const checkOut = async () => {
  try {
    const userStr = localStorage.getItem('user')
    if (!userStr) {
      alert('로그인 정보가 없습니다. 다시 로그인 해주세요.')
      return
    }
    
    const user = JSON.parse(userStr)
    const user_id = user.id

    // 현재 시간을 ISO 문자열로 전송
    const now = new Date()
    const response = await api.post('/api/attendance/checkout', {
      user_id: parseInt(user_id),
      date: now.toISOString().split('T')[0],
      check_out_time: now.toISOString()
    })

    if (response.data.success) {
      // 출퇴근 기록 다시 로드하여 최신 데이터 반영
      await loadTodayAttendance()
      await loadAttendanceRecords()
      
      alert(`퇴근이 기록되었습니다: ${response.data.check_out_time}`)
      
      // 조퇴 체크
      const checkOutTime = parseTimeToMinutes(attendanceState.todayAttendance.checkOutTime)
      const standardEndTime = 18 * 60
      
      if (checkOutTime < standardEndTime) {
        alert(`조퇴입니다. 표준 퇴근 시간은 18:00입니다. (${Math.floor((standardEndTime - checkOutTime) / 60)}시간 ${(standardEndTime - checkOutTime) % 60}분 일찍 퇴근)`)
      }
    } else {
      alert('퇴근 기록에 실패했습니다.')
    }
  } catch (error) {
    console.error('[checkOut] 오류:', error)
    if (error.response?.status === 404) {
      alert('출근 기록이 없습니다. 먼저 출근해주세요.')
    } else if (error.response?.status === 409) {
      alert('이미 퇴근 기록이 있습니다.')
    } else {
      alert('퇴근 기록 중 오류가 발생했습니다.')
    }
  }
}

const loadAttendanceByDate = async (date) => {
  try {
    const userStr = localStorage.getItem('user')
    if (!userStr) {
      console.log('[loadAttendanceByDate] 로그인 정보 없음')
      return
    }
    
    const user = JSON.parse(userStr)
    const user_id = user.id

    console.log('[loadAttendanceByDate] 요청 URL:', `/api/attendance/date/${user_id}/${date}`)
    const response = await api.get(`/api/attendance/date/${user_id}/${date}`)
    
    console.log('[loadAttendanceByDate] 응답 데이터:', response.data)
    
    if (response.data && response.data.work_date) {
      // UTC 시간을 로컬 시간으로 변환
      const formatTime = (timeString) => {
        if (!timeString) return null
        const date = new Date(timeString)
        return date.toLocaleTimeString('ko-KR', { 
          hour: '2-digit', 
          minute: '2-digit',
          hour12: false 
        })
      }

      attendanceState.todayAttendance = {
        date: response.data.work_date,
        checkInTime: formatTime(response.data.check_in_time),
        checkOutTime: formatTime(response.data.check_out_time),
        originalCheckInTime: response.data.check_in_time,
        originalCheckOutTime: response.data.check_out_time,
        workTime: response.data.check_in_time && response.data.check_out_time ? 
          formatMinutesToHours((new Date(response.data.check_out_time) - new Date(response.data.check_in_time)) / 1000 / 60) : '',
        status: response.data.status
      }
      console.log('[loadAttendanceByDate] todayAttendance 업데이트됨:', attendanceState.todayAttendance)
    } else {
      // 해당 날짜에 기록이 없는 경우
      attendanceState.todayAttendance = {
        date: date,
        checkInTime: null,
        checkOutTime: null,
        originalCheckInTime: null,
        originalCheckOutTime: null,
        workTime: '',
        status: null
      }
      console.log('[loadAttendanceByDate] 해당 날짜에 기록 없음, 초기화됨')
    }
  } catch (error) {
    console.error('[loadAttendanceByDate] 오류:', error)
  }
}

const loadTodayAttendance = async () => {
  const today = new Date().toISOString().split('T')[0]
  await loadAttendanceByDate(today)
}

const loadAttendanceRecords = async () => {
  try {
    const userStr = localStorage.getItem('user')
    if (!userStr) {
      console.log('[loadAttendanceRecords] 로그인 정보 없음')
      return
    }
    
    const user = JSON.parse(userStr)
    const user_id = user.id

    console.log('[loadAttendanceRecords] 요청 URL:', `/api/attendance/records/${user_id}`)
    console.log('[loadAttendanceRecords] user_id:', user_id, '타입:', typeof user_id)
    const response = await api.get(`/api/attendance/records/${user_id}`)
    
    if (response.data) {
      // UTC 시간을 로컬 시간으로 변환
      const formatTime = (timeString) => {
        if (!timeString) return null
        const date = new Date(timeString)
        return date.toLocaleTimeString('ko-KR', { 
          hour: '2-digit', 
          minute: '2-digit',
          hour12: false 
        })
      }

      attendanceState.attendanceRecords = response.data.map(record => ({
        date: record.work_date,
        checkInTime: formatTime(record.check_in_time),
        checkOutTime: formatTime(record.check_out_time),
        workHours: record.check_in_time && record.check_out_time ? 
          (new Date(record.check_out_time) - new Date(record.check_in_time)) / 1000 / 60 / 60 : null,
        status: record.status
      }))
    }
  } catch (error) {
    console.error('[loadAttendanceRecords] 오류:', error)
  }
}

const calculateWeeklyWorkTime = () => {
  const today = new Date()
  const dayOfWeek = today.getDay()
  const startOfWeek = new Date(today)
  startOfWeek.setDate(today.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1))
  
  const endOfWeek = new Date(startOfWeek)
  endOfWeek.setDate(startOfWeek.getDate() + 4)
  
  const formatDate = (date) => {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }
  
  const startDate = formatDate(startOfWeek)
  const endDate = formatDate(endOfWeek)
  
  const weekRecords = attendanceState.attendanceRecords.filter(record => {
    return record.date >= startDate && record.date <= endDate
  })
  
  let totalMinutes = 0
  
  weekRecords.forEach(record => {
    if (record.checkInTime && record.checkOutTime) {
      const checkIn = parseTimeToMinutes(record.checkInTime)
      const checkOut = parseTimeToMinutes(record.checkOutTime)
      totalMinutes += (checkOut - checkIn)
    }
  })
  
  attendanceState.weeklyWorkMinutes = totalMinutes
}

const parseTimeToMinutes = (timeString) => {
  const [hours, minutes] = timeString.split(':').map(Number)
  return hours * 60 + minutes
}

const formatMinutesToHours = (minutes) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}시간 ${mins}분`
}

const startTimer = () => {
  updateCurrentTime()
  attendanceState.timerInterval = setInterval(updateCurrentTime, 1000)
}

const stopTimer = () => {
  if (attendanceState.timerInterval) {
    clearInterval(attendanceState.timerInterval)
    attendanceState.timerInterval = null
  }
}

const initializeAttendance = async () => {
  startTimer()
  await loadTodayAttendance()
  await loadAttendanceRecords()
  calculateWeeklyWorkTime()
}

// Pinia와 동일한 인터페이스를 제공하는 함수
export const useAttendanceStore = () => {
  return {
    // State
    currentTime: computed(() => attendanceState.currentTime),
    todayAttendance: computed(() => attendanceState.todayAttendance),
    attendanceRecords: computed(() => attendanceState.attendanceRecords),
    workType: computed(() => attendanceState.workType),
    weeklyWorkMinutes: computed(() => attendanceState.weeklyWorkMinutes),
    weeklyWorkGoal: computed(() => attendanceState.weeklyWorkGoal),
    
    // Getters
    workingHours,
    weeklyWorkPercentage,
    getWorkStatus,
    getWorkStatusClass,
    
    // Actions
    updateCurrentTime,
    checkIn,
    checkOut,
    loadAttendanceByDate,
    loadTodayAttendance,
    loadAttendanceRecords,
    calculateWeeklyWorkTime,
    parseTimeToMinutes,
    formatMinutesToHours,
    startTimer,
    stopTimer,
    initializeAttendance
  }
} 