// 날짜 포맷팅
export const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ko-KR', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric', 
    weekday: 'long' 
  });
};

// 시간 포맷팅 (문자열에서 시간만 추출)
export const formatTime = (dateTimeString) => {
  if (!dateTimeString) return '';
  // "2023-07-24T15:00" 또는 "2023-07-24 15:00"
  const timePart = dateTimeString.split('T')[1] || dateTimeString.split(' ')[1];
  return timePart ? timePart.slice(0, 5) : '';
};

// 분을 시간:분 형식으로 변환
export const formatMinutesToHours = (minutes) => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `${hours}시간 ${mins}분`;
};

// 시간을 분으로 변환
export const parseTimeToMinutes = (timeString) => {
  const [hours, minutes] = timeString.split(':').map(Number);
  return hours * 60 + minutes;
};

// 분을 시간 형식으로 변환
export const formatMinuteToTime = (minute) => {
  const h = Math.floor(minute / 60);
  const m = minute % 60;
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
}; 