import { ref, computed } from 'vue';
import api from '@/plugins/axios';

export function useAttendance() {
  const currentTime = ref('');
  const timerInterval = ref(null);
  
  const attendanceRecords = ref([]);

  const workType = ref('fixed'); // 'fixed': 고정형(9to6), 'flexible': 자율형(주 40시간)
  const weeklyWorkMinutes = ref(0);
  const weeklyWorkGoal = 2400; // 주 40시간 = 2400분

  const todayAttendance = ref({
    date: new Date().toISOString().split('T')[0],
    checkInTime: null,
    checkOutTime: null,
    workTime: ''
  });

  // 근무 시간 계산
  const workingHours = computed(() => {
    if (!todayAttendance.value.checkInTime || !todayAttendance.value.checkOutTime) {
      return '-';
    }
    
    const checkIn = parseTimeToMinutes(todayAttendance.value.checkInTime);
    const checkOut = parseTimeToMinutes(todayAttendance.value.checkOutTime);
    const workMinutes = checkOut - checkIn;
    
    const hours = Math.floor(workMinutes / 60);
    const minutes = workMinutes % 60;
    
    return `${hours}시간 ${minutes}분`;
  });

  // 주간 근무 시간 퍼센트
  const weeklyWorkPercentage = computed(() => {
    return Math.min(100, (weeklyWorkMinutes.value / weeklyWorkGoal) * 100);
  });

  // 근무 상태 계산
  const getWorkStatus = computed(() => {
    if (!todayAttendance.value.checkInTime) {
      return '';
    }
    
    const checkInTime = parseTimeToMinutes(todayAttendance.value.checkInTime);
    const standardStartTime = 9 * 60;
    
    if (todayAttendance.value.checkOutTime) {
      const checkOutTime = parseTimeToMinutes(todayAttendance.value.checkOutTime);
      const standardEndTime = 18 * 60;
      
      if (checkInTime > standardStartTime && checkOutTime < standardEndTime) {
        return '지각 및 조퇴';
      } else if (checkInTime > standardStartTime) {
        return '지각';
      } else if (checkOutTime < standardEndTime) {
        return '조퇴';
      } else {
        return '정상 출퇴근';
      }
    } else {
      return checkInTime > standardStartTime ? '지각' : '근무중';
    }
  });

  // 근무 상태에 따른 클래스
  const getWorkStatusClass = computed(() => {
    const status = getWorkStatus.value;
    if (status.includes('지각') || status.includes('조퇴')) {
      return 'status-warning';
    } else if (status === '정상 출퇴근') {
      return 'status-success';
    }
    return '';
  });

  // 현재 시간 업데이트
  const updateCurrentTime = () => {
    const now = new Date();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const seconds = String(now.getSeconds()).padStart(2, '0');
    currentTime.value = `${hours}:${minutes}:${seconds}`;
  };

  // 출근
  const checkIn = async () => {
    try {
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        alert('로그인 정보가 없습니다. 다시 로그인 해주세요.');
        return;
      }
      
      const user = JSON.parse(userStr);
      const user_id = user.id;

      // 현재 시간을 ISO 문자열로 전송
      const now = new Date();
      const response = await api.post('/api/attendance/checkin', {
        user_id: parseInt(user_id),
        date: now.toISOString().split('T')[0],
        check_in_time: now.toISOString()
      });

      if (response.data.success) {
        // 출퇴근 기록 다시 로드하여 최신 데이터 반영
        await loadTodayAttendance();
        await loadAttendanceRecords();
        
        alert(`출근이 기록되었습니다: ${response.data.check_in_time}`);
        
        // 지각 체크
        const checkInTime = parseTimeToMinutes(todayAttendance.value.checkInTime);
        const standardStartTime = 9 * 60;
        
        if (checkInTime > standardStartTime) {
          alert(`지각입니다. 표준 출근 시간은 09:00입니다. (${Math.floor((checkInTime - standardStartTime) / 60)}시간 ${(checkInTime - standardStartTime) % 60}분 지각)`);
        }
      } else {
        alert('출근 기록에 실패했습니다.');
      }
    } catch (error) {
      console.error('[checkIn] 오류:', error);
      if (error.response?.status === 409) {
        alert('이미 출근 기록이 있습니다.');
      } else {
        alert('출근 기록 중 오류가 발생했습니다.');
      }
    }
  };

  // 퇴근
  const checkOut = async () => {
    try {
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        alert('로그인 정보가 없습니다. 다시 로그인 해주세요.');
        return;
      }
      
      const user = JSON.parse(userStr);
      const user_id = user.id;

      // 현재 시간을 ISO 문자열로 전송
      const now = new Date();
      const response = await api.post('/api/attendance/checkout', {
        user_id: parseInt(user_id),
        date: now.toISOString().split('T')[0],
        check_out_time: now.toISOString()
      });

      if (response.data.success) {
        // 출퇴근 기록 다시 로드하여 최신 데이터 반영
        await loadTodayAttendance();
        await loadAttendanceRecords();
        
        alert(`퇴근이 기록되었습니다: ${response.data.check_out_time}`);
        
        // 조퇴 체크
        const checkOutTime = parseTimeToMinutes(todayAttendance.value.checkOutTime);
        const standardEndTime = 18 * 60;
        
        if (checkOutTime < standardEndTime) {
          alert(`조퇴입니다. 표준 퇴근 시간은 18:00입니다. (${Math.floor((standardEndTime - checkOutTime) / 60)}시간 ${(standardEndTime - checkOutTime) % 60}분 일찍 퇴근)`);
        }
      } else {
        alert('퇴근 기록에 실패했습니다.');
      }
    } catch (error) {
      console.error('[checkOut] 오류:', error);
      if (error.response?.status === 404) {
        alert('출근 기록이 없습니다. 먼저 출근해주세요.');
      } else if (error.response?.status === 409) {
        alert('이미 퇴근 기록이 있습니다.');
      } else {
        alert('퇴근 기록 중 오류가 발생했습니다.');
      }
    }
  };

  // 출퇴근 기록 저장
  const saveAttendance = () => {
    console.log('출퇴근 기록 저장:', todayAttendance.value);
    
    const index = attendanceRecords.value.findIndex(record => record.date === todayAttendance.value.date);
    if (index !== -1) {
      attendanceRecords.value.splice(index, 1, { ...todayAttendance.value });
    } else {
      attendanceRecords.value.push({ ...todayAttendance.value });
    }
  };

  // 특정 날짜 출퇴근 기록 불러오기
  const loadAttendanceByDate = async (date) => {
    try {
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        console.log('[loadAttendanceByDate] 로그인 정보 없음');
        return;
      }
      
      const user = JSON.parse(userStr);
      const user_id = user.id;

      console.log('[loadAttendanceByDate] 요청 URL:', `/api/attendance/date/${user_id}/${date}`);
      const response = await api.get(`/api/attendance/date/${user_id}/${date}`);
      
      console.log('[loadAttendanceByDate] 응답 데이터:', response.data);
      
      if (response.data && response.data.work_date) {
        // UTC 시간을 로컬 시간으로 변환
        const formatTime = (timeString) => {
          if (!timeString) return null;
          const date = new Date(timeString);
          return date.toLocaleTimeString('ko-KR', { 
            hour: '2-digit', 
            minute: '2-digit',
            hour12: false 
          });
        };

        todayAttendance.value = {
          date: response.data.work_date,
          checkInTime: formatTime(response.data.check_in_time),
          checkOutTime: formatTime(response.data.check_out_time),
          workTime: response.data.check_in_time && response.data.check_out_time ? 
            formatMinutesToHours((new Date(response.data.check_out_time) - new Date(response.data.check_in_time)) / 1000 / 60) : '',
          status: response.data.status
        };
        console.log('[loadAttendanceByDate] todayAttendance 업데이트됨:', todayAttendance.value);
      } else {
        // 해당 날짜에 기록이 없는 경우
        todayAttendance.value = {
          date: date,
          checkInTime: null,
          checkOutTime: null,
          workTime: '',
          status: null
        };
        console.log('[loadAttendanceByDate] 해당 날짜에 기록 없음, 초기화됨');
      }
    } catch (error) {
      console.error('[loadAttendanceByDate] 오류:', error);
    }
  };

  // 오늘 출퇴근 기록 불러오기
  const loadTodayAttendance = async () => {
    const today = new Date().toISOString().split('T')[0];
    await loadAttendanceByDate(today);
  };

  // 출퇴근 기록 목록 불러오기
  const loadAttendanceRecords = async () => {
    try {
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        console.log('[loadAttendanceRecords] 로그인 정보 없음');
        return;
      }
      
      const user = JSON.parse(userStr);
      const user_id = user.id;

      console.log('[loadAttendanceRecords] 요청 URL:', `/api/attendance/records/${user_id}`);
      console.log('[loadAttendanceRecords] user_id:', user_id, '타입:', typeof user_id);
      const response = await api.get(`/api/attendance/records/${user_id}`);
      
      if (response.data) {
        // UTC 시간을 로컬 시간으로 변환
        const formatTime = (timeString) => {
          if (!timeString) return null;
          const date = new Date(timeString);
          return date.toLocaleTimeString('ko-KR', { 
            hour: '2-digit', 
            minute: '2-digit',
            hour12: false 
          });
        };

        attendanceRecords.value = response.data.map(record => ({
          date: record.work_date,
          checkInTime: formatTime(record.check_in_time),
          checkOutTime: formatTime(record.check_out_time),
          workHours: record.check_in_time && record.check_out_time ? 
            (new Date(record.check_out_time) - new Date(record.check_in_time)) / 1000 / 60 / 60 : null,
          status: record.status
        }));
      }
    } catch (error) {
      console.error('[loadAttendanceRecords] 오류:', error);
    }
  };

  // 주간 근무 시간 계산
  const calculateWeeklyWorkTime = () => {
    const today = new Date();
    const dayOfWeek = today.getDay();
    const startOfWeek = new Date(today);
    startOfWeek.setDate(today.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1));
    
    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 4);
    
    const formatDate = (date) => {
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
    };
    
    const startDate = formatDate(startOfWeek);
    const endDate = formatDate(endOfWeek);
    
    const weekRecords = attendanceRecords.value.filter(record => {
      return record.date >= startDate && record.date <= endDate;
    });
    
    let totalMinutes = 0;
    
    weekRecords.forEach(record => {
      if (record.checkInTime && record.checkOutTime) {
        const checkIn = parseTimeToMinutes(record.checkInTime);
        const checkOut = parseTimeToMinutes(record.checkOutTime);
        totalMinutes += (checkOut - checkIn);
      }
    });
    
    weeklyWorkMinutes.value = totalMinutes;
  };

  // 시간을 분으로 변환
  const parseTimeToMinutes = (timeString) => {
    const [hours, minutes] = timeString.split(':').map(Number);
    return hours * 60 + minutes;
  };

  // 분을 시간:분 형식으로 변환
  const formatMinutesToHours = (minutes) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}시간 ${mins}분`;
  };

  // 타이머 시작
  const startTimer = () => {
    updateCurrentTime();
    timerInterval.value = setInterval(updateCurrentTime, 1000);
  };

  // 타이머 정리
  const stopTimer = () => {
    if (timerInterval.value) {
      clearInterval(timerInterval.value);
      timerInterval.value = null;
    }
  };

  // 초기화
  const initializeAttendance = async () => {
    startTimer();
    await loadTodayAttendance();
    await loadAttendanceRecords();
    calculateWeeklyWorkTime();
  };

  return {
    currentTime,
    attendanceRecords,
    workType,
    weeklyWorkMinutes,
    weeklyWorkGoal,
    todayAttendance,
    workingHours,
    weeklyWorkPercentage,
    getWorkStatus,
    getWorkStatusClass,
    updateCurrentTime,
    checkIn,
    checkOut,
    saveAttendance,
    loadTodayAttendance,
    loadAttendanceByDate,
    loadAttendanceRecords,
    calculateWeeklyWorkTime,
    parseTimeToMinutes,
    formatMinutesToHours,
    startTimer,
    stopTimer,
    initializeAttendance
  };
} 