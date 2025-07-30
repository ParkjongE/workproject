import { ref, computed } from 'vue';

export function useCalendar() {
  const today = new Date();
  const selectedDate = ref(`${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`);
  const currentYear = ref(today.getFullYear());
  const currentMonth = ref(today.getMonth());
  const weekdays = ['일', '월', '화', '수', '목', '금', '토'];
  const showYearDropdown = ref(false);
  const showMonthDropdown = ref(false);

  // 선택 가능한 년도 목록 (현재 년도 기준 전후 5년)
  const availableYears = computed(() => {
    const currentYearValue = new Date().getFullYear();
    const years = [];
    for (let i = currentYearValue - 5; i <= currentYearValue + 5; i++) {
      years.push(i);
    }
    return years;
  });

  // 달력 날짜들 계산
  const calendarDays = computed(() => {
    const days = [];
    const firstDay = new Date(currentYear.value, currentMonth.value, 1);
    const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0);
    
    // 이전 달의 마지막 날짜들
    const prevMonthLastDate = new Date(currentYear.value, currentMonth.value, 0).getDate();
    const firstDayOfWeek = firstDay.getDay();
    
    for (let i = firstDayOfWeek - 1; i >= 0; i--) {
      const prevMonthYear = currentMonth.value === 0 ? currentYear.value - 1 : currentYear.value;
      const prevMonth = currentMonth.value === 0 ? 11 : currentMonth.value - 1;
      const day = prevMonthLastDate - i;
      days.push({
        date: `${prevMonthYear}-${String(prevMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`,
        dayNumber: day,
        isCurrentMonth: false
      });
    }
    
    // 현재 달의 날짜들
    for (let i = 1; i <= lastDay.getDate(); i++) {
      days.push({
        date: `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`,
        dayNumber: i,
        isCurrentMonth: true
      });
    }
    
    // 다음 달의 시작 날짜들
    const lastDayOfWeek = lastDay.getDay();
    for (let i = 1; i < 7 - lastDayOfWeek; i++) {
      const nextMonthYear = currentMonth.value === 11 ? currentYear.value + 1 : currentYear.value;
      const nextMonth = currentMonth.value === 11 ? 0 : currentMonth.value + 1;
      days.push({
        date: `${nextMonthYear}-${String(nextMonth + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`,
        dayNumber: i,
        isCurrentMonth: false
      });
    }
    
    return days;
  });

  // 이전 달
  const prevMonth = () => {
    if (currentMonth.value === 0) {
      currentYear.value--;
      currentMonth.value = 11;
    } else {
      currentMonth.value--;
    }
  };

  // 다음 달
  const nextMonth = () => {
    if (currentMonth.value === 11) {
      currentYear.value++;
      currentMonth.value = 0;
    } else {
      currentMonth.value++;
    }
  };

  // 날짜 선택
  const selectDay = (date) => {
    selectedDate.value = date;
  };

  // 선택된 날짜인지 확인
  const isSelectedDay = (date) => {
    return selectedDate.value === date;
  };

  // 년도 드롭다운 토글
  const toggleYearDropdown = () => {
    showYearDropdown.value = !showYearDropdown.value;
    if (showYearDropdown.value) {
      showMonthDropdown.value = false;
    }
  };

  // 월 드롭다운 토글
  const toggleMonthDropdown = () => {
    showMonthDropdown.value = !showMonthDropdown.value;
    if (showMonthDropdown.value) {
      showYearDropdown.value = false;
    }
  };

  // 년도 선택
  const selectYear = (year) => {
    currentYear.value = year;
    showYearDropdown.value = false;
  };

  // 월 선택
  const selectMonth = (month) => {
    currentMonth.value = month;
    showMonthDropdown.value = false;
  };

  // 초기화
  const initializeCalendar = () => {
    const today = new Date();
    selectedDate.value = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`;
  };

  return {
    currentYear,
    currentMonth,
    selectedDate,
    weekdays,
    showYearDropdown,
    showMonthDropdown,
    availableYears,
    calendarDays,
    prevMonth,
    nextMonth,
    selectDay,
    isSelectedDay,
    toggleYearDropdown,
    toggleMonthDropdown,
    selectYear,
    selectMonth,
    initializeCalendar
  };
} 