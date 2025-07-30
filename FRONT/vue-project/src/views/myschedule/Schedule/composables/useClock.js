import { ref } from 'vue';

export function useClock() {
  const hoveredIdx = ref(null);
  const showQuickAddForm = ref(false);
  const colorAnimationActive = ref(false);

  const quickAddForm = ref({
    title: '',
    description: '',
    startTime: '',
    endTime: '',
    color: '#4CAF50',
    ampm: '',
  });

  const colorOptions = [
    '#4CAF50', // 초록
    '#2196F3', // 파랑
    '#FF9800', // 주황
    '#F44336', // 빨강
    '#9C27B0', // 보라
    '#00BCD4', // 청록
    '#FFEB3B', // 노랑
    '#795548'  // 갈색
  ];

  // 시계 숫자 스타일 계산
  const getHourNumberStyle = (hour) => {
    const clockSize = 520;
    const radius = 240;
    const centerX = clockSize / 2;
    const centerY = clockSize / 2;
    const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
    const x = Math.cos(angle) * radius + centerX;
    const y = Math.sin(angle) * radius + centerY;
    return {
      position: 'absolute',
      left: `${x}px`,
      top: `${y}px`,
      transform: 'translate(-50%, -50%)',
      fontSize: '28px',
      fontWeight: 'bold',
      color: '#333',
      zIndex: 2
    };
  };

  // 일정 세그먼트 스타일 계산
  const getScheduleSegmentStyle = (schedule) => {
    const startParts = schedule.start.split('T')[1].split(':');
    const endParts = schedule.end.split('T')[1].split(':');
    let startHour = parseInt(startParts[0]);
    const startMinute = parseInt(startParts[1]);
    let endHour = parseInt(endParts[0]);
    const endMinute = parseInt(endParts[1]);
    
    function to12Hour(hour) {
      if (hour === 0) return 12;
      if (hour > 12) return hour - 12;
      return hour;
    }
    
    const startHour12 = to12Hour(startHour);
    const endHour12 = to12Hour(endHour);
    const startDeg = ((startHour12 % 12) * 30 + startMinute * 0.5);
    let endDeg = ((endHour12 % 12) * 30 + endMinute * 0.5);
    if (endDeg <= startDeg) endDeg += 360;
    const sweepDeg = endDeg - startDeg;
    const clockSize = 520;
    
    return {
      position: 'absolute',
      left: '50%',
      top: '50%',
      width: `${clockSize}px`,
      height: `${clockSize}px`,
      transform: 'translate(-50%, -50%)',
      background: `conic-gradient(from ${startDeg}deg, ${schedule.color} 0deg, ${schedule.color} ${sweepDeg}deg, transparent ${sweepDeg}deg 360deg)`,
      borderRadius: '50%',
      opacity: '0.7',
      zIndex: 0,
      pointerEvents: 'none',
      cursor: 'pointer'
    };
  };

  // 시간대 선택 슬롯 스타일
  const getTimeSlotStyle = (hour) => {
    const clockSize = 520;
    const radius = 240;
    const centerX = clockSize / 2;
    const centerY = clockSize / 2;
    const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
    const x = Math.cos(angle) * radius + centerX;
    const y = Math.sin(angle) * radius + centerY;
    
    return {
      position: 'absolute',
      left: `${x}px`,
      top: `${y}px`,
      transform: 'translate(-50%, -50%)',
      zIndex: 4
    };
  };

  // 빠른 추가 폼 위치 계산
  const getQuickAddFormStyle = () => {
    return {
      position: 'absolute',
      top: '50%',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      zIndex: 1000
    };
  };

  // 시계 클릭 이벤트 처리
  const onClockFaceClick = (event) => {
    if (showQuickAddForm.value) return;
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const x = event.clientX - centerX;
    const y = event.clientY - centerY;
    let angle = Math.atan2(y, x) * 180 / Math.PI;
    angle = (angle + 450) % 360;
    let hour = Math.floor(angle / 30);
    hour = (hour + 12) % 12;
    if (hour === 0) hour = 12;
    // 수정: 클릭한 시각을 그대로 사용 (am/pm 변환 및 +12 보정 제거)
    let realHour = hour;
    showQuickAddForm.value = true;
    // 기존 description은 유지하고 나머지만 초기화
    const currentDescription = quickAddForm.value.description || '';
    quickAddForm.value = {
      title: '',
      description: currentDescription,
      startTime: `${String(realHour).padStart(2, '0')}:00`,
      endTime: `${String((realHour % 24 + 1)).padStart(2, '0')}:00`,
      color: '#4CAF50',
      ampm: '',
    };
  };

  // 색상 선택
  const selectColor = (color) => {
    quickAddForm.value.color = color;
    colorAnimationActive.value = true;
    setTimeout(() => {
      colorAnimationActive.value = false;
    }, 300);
  };

  // 빠른 일정 추가 폼 닫기
  const closeQuickAddForm = () => {
    showQuickAddForm.value = false;
    quickAddForm.value = {
      title: '',
      description: '',
      startTime: '',
      endTime: '',
      color: '#4CAF50',
      ampm: '',
    };
  };

  // 시간대 선택 확인
  const getTimeSlotSelected = (hour) => {
    if (showQuickAddForm.value && quickAddForm.value.startTime && quickAddForm.value.endTime) {
      const start = parseTimeToMinutes(quickAddForm.value.startTime);
      const end = parseTimeToMinutes(quickAddForm.value.endTime);
      return hour >= start && hour <= end;
    }
    return false;
  };

  // 시간을 분으로 변환
  const parseTimeToMinutes = (timeString) => {
    const [hours, minutes] = timeString.split(':').map(Number);
    return hours * 60 + minutes;
  };

  // 시계 외부 클릭 처리
  const handleOutsideClick = (event) => {
    const clockFace = event.target.closest('.clock-face');
    const quickAddFormElement = event.target.closest('.quick-add-form');
    
    if (!clockFace && !quickAddFormElement && showQuickAddForm.value) {
      closeQuickAddForm();
    }
  };

  return {
    hoveredIdx,
    showQuickAddForm,
    quickAddForm,
    colorOptions,
    colorAnimationActive,
    getHourNumberStyle,
    getScheduleSegmentStyle,
    getTimeSlotStyle,
    getQuickAddFormStyle,
    onClockFaceClick,
    selectColor,
    closeQuickAddForm,
    getTimeSlotSelected,
    parseTimeToMinutes,
    handleOutsideClick
  };
} 