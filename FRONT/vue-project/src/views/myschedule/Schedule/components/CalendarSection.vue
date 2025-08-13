<template>
  <CalendarSectionStyle>
    <CalendarHeader>
      <button @click="prevMonth">&lt;</button>
      <div class="date-selector">
        <DropdownContainer>
          <DropdownToggle @click="toggleYearDropdown">
            <span>{{ currentYear }}년</span>
            <DropdownArrow :class="{ 'open': showYearDropdown }" />
          </DropdownToggle>
          <DropdownMenu v-if="showYearDropdown">
            <DropdownItem
              v-for="year in availableYears"
              :key="year"
              :class="{ 'selected': year === currentYear }"
              @click="selectYear(year)"
            >
              {{ year }}년
            </DropdownItem>
          </DropdownMenu>
        </DropdownContainer>
        <DropdownContainer>
          <DropdownToggle @click="toggleMonthDropdown">
            <span>{{ currentMonth + 1 }}월</span>
            <DropdownArrow :class="{ 'open': showMonthDropdown }" />
          </DropdownToggle>
          <DropdownMenu v-if="showMonthDropdown">
            <DropdownItem
              v-for="month in 12"
              :key="month"
              :class="{ 'selected': month === currentMonth + 1 }"
              @click="selectMonth(month - 1)"
            >
              {{ month }}월
            </DropdownItem>
          </DropdownMenu>
        </DropdownContainer>
      </div>
      <button @click="nextMonth">&gt;</button>
    </CalendarHeader>
    <CalendarGrid>
      <div class="weekday-header" v-for="day in weekdays" :key="day">{{ day }}</div>
      <CalendarDay
        v-for="day in calendarDays"
        :key="day.date"
        :selected="isSelectedDay(day.date)"
        :otherMonth="!day.isCurrentMonth"
        @click="selectDay(day.date)"
      >
        <DayNumber>{{ day.dayNumber }}</DayNumber>
        <DaySchedules>
          <DayScheduleItem
            v-for="(schedule, idx) in getDaySchedules(day.date) || []"
            :key="idx"
            :color="schedule.color"
            @click.stop="selectSchedule(schedule)"
          >
            {{ schedule.title }}
          </DayScheduleItem>
        </DaySchedules>
      </CalendarDay>
    </CalendarGrid>
  </CalendarSectionStyle>
</template>

<script>
import * as styled from '../Schedule.styled.js';
import { watch } from 'vue';
// Composables
import { useCalendar } from '../composables/useCalendar';

export default {
  name: 'CalendarSection',
  components: {
    ...styled,
  },
  props: {
    mySchedules: {
      type: Array,
      default: () => []
    }
  },
  emits: ['day-selected', 'schedule-selected'],
  setup(props, { emit }) {
    // Composables 사용
    const {
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
    } = useCalendar();

    // props로 받은 mySchedules를 사용하여 getDaySchedules 함수 생성
    const getDaySchedules = (date) => {
      if (!date) return [];
      // props.mySchedules에 직접 접근하여 반응성 보장
      const schedules = props.mySchedules || [];
      return schedules.filter(schedule => {
        const scheduleDate = schedule.start.split('T')[0];
        return scheduleDate === date;
      });
    };

    const selectSchedule = (schedule) => {
      // 일정 선택 로직
      emit('schedule-selected', schedule);
    };

    // 날짜 선택 핸들러
    const handleSelectDay = (date) => {
      selectDay(date);
      emit('day-selected', date);
    };

    // 일정 선택 핸들러
    const handleSelectSchedule = (schedule) => {
      selectSchedule(schedule);
      emit('schedule-selected', schedule);
    };

    // selectedDate가 바뀌면(초기 진입 포함) 해당 날짜로 selectDay 호출
    watch(
      () => props.selectedDate,
      (newDate) => {
        if (newDate) selectDay(newDate);
      },
      { immediate: true }
    );

    return {
      // Calendar
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
      selectDay: handleSelectDay,
      isSelectedDay,
      toggleYearDropdown,
      toggleMonthDropdown,
      selectYear,
      selectMonth,
      initializeCalendar,
      
      // Schedule
      getDaySchedules,
      selectSchedule: handleSelectSchedule
    };
  },
  mounted() {
    this.initializeCalendar();
  }
};
</script> 