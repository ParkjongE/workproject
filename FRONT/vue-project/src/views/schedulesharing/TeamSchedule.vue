<template>
  <TeamScheduleContainer>
    <SideBar />
    <ContentArea>
      <h1>팀 일정 공유</h1>
      <DashboardContainer>
        <!-- 팀 목록 그리드 -->
        <TeamGrid>
          <TeamCard 
            v-for="(team, index) in teams" 
            :key="index"
          >
            <TeamTitle>{{ team.name }}</TeamTitle>
            <TaskList>
              <li v-for="(task, taskIndex) in team.tasks.slice(0, 3)" :key="taskIndex">
                <span>{{ task.title }}</span>
                <Deadline>{{ task.deadline }}</Deadline>
              </li>
            </TaskList>
            <More>
              {{ team.tasks.length > 3 ? '외 ' + (team.tasks.length - 3) + '건 더 있음' : '\u00A0' }}
            </More>
            <TeamMembers>
              팀원: 
              <template v-for="(member, memberIndex) in team.members" :key="memberIndex">
                <TeamLeader 
                  v-if="member === team.leader"
                >
                  {{ member }}(팀장){{ memberIndex < team.members.length - 1 ? ', ' : '' }}
                </TeamLeader>
                <span v-else>
                  {{ member }}{{ memberIndex < team.members.length - 1 ? ', ' : '' }}
                </span>
              </template>
            </TeamMembers>
          </TeamCard>
        </TeamGrid>
      </DashboardContainer>
    </ContentArea>
  </TeamScheduleContainer>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  TeamScheduleContainer, ContentArea, DashboardContainer, TeamGrid, TeamCard, TeamTitle, TaskList, Deadline, More, TeamMembers, TeamLeader
} from './TeamSchedule.styled.js';

export default {
  name: 'TeamSchedule',
  components: {
    SideBar,
    TeamScheduleContainer, ContentArea, DashboardContainer, TeamGrid, TeamCard, TeamTitle, TaskList, Deadline, More, TeamMembers, TeamLeader
  },
  data() {
    return {
      teams: [
        {
          name: '개발1팀',
          leader: '홍길동',
          members: ['홍길동', '박주호', '이민수', '최수진'],
          tasks: [
            { title: '신규 서비스 런칭', deadline: '2024-08-10' },
            { title: 'API 리팩토링', deadline: '2024-08-20' },
            { title: '버그 수정', deadline: '2024-08-25' }
          ]
        },
        {
          name: '개발2팀',
          leader: '김지훈',
          members: ['김지훈', '정유진', '오세훈', '한지민'],
          tasks: [
            { title: '앱 디자인 개선', deadline: '2024-08-12' },
            { title: '로그인 기능 개선', deadline: '2024-08-18' }
          ]
        },
        {
          name: '영업1팀',
          leader: '이수진',
          members: ['이수진', '박성호', '최민지'],
          tasks: [
            { title: '고객사 미팅', deadline: '2024-08-09' },
            { title: '신규 계약 추진', deadline: '2024-08-22' }
          ]
        },
        {
          name: '영업2팀',
          leader: '정다은',
          members: ['정다은', '이준호'],
          tasks: [
            { title: '시장조사', deadline: '2024-08-15' }
          ]
        }
      ],
      memberSchedules: {
        '홍길동': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-13'], leave: ['2024-08-12'] },
        '박주호': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '이민수': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '최수진': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '김지훈': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '정유진': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '오세훈': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '한지민': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '이수진': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '박성호': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '최민지': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '정다은': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] },
        '이준호': { work: ['2024-08-01','2024-08-02','2024-08-05','2024-08-06','2024-08-07','2024-08-08','2024-08-09','2024-08-12','2024-08-13'], leave: [] }
      }
    };
  },
  methods: {
    isOnLeave(member) {
      return this.memberSchedules[member]?.leave.includes('2024-08-12');
    },
    makeCalendar(year, month, member, mini) {
      const firstDay = new Date(year, month-1, 1).getDay();
      const lastDate = new Date(year, month, 0).getDate();
      let html = `<table class="calendar-table ${mini ? 'mini' : ''}">`;
      html += '<tr class="calendar-header"><th>일</th><th>월</th><th>화</th><th>수</th><th>목</th><th>금</th><th>토</th></tr><tr>';
      
      for(let i=0; i<firstDay; i++) html += '<td></td>';
      
      for(let d=1; d<=lastDate; d++) {
        const dateStr = `${year}-08-${d.toString().padStart(2,'0')}`;
        let cellClass = '';
        
        if(this.memberSchedules[member]?.leave.includes(dateStr)) {
          cellClass = 'leave-day';
        } else if(this.memberSchedules[member]?.work.includes(dateStr)) {
          cellClass = 'work-day';
        }
        
        html += `<td class="${cellClass}">${d}</td>`;
        
        if((firstDay+d)%7===0) html += '</tr><tr>';
      }
      
      html += '</tr></table>';
      return html;
    },
    clockSVG(isLeave, mini) {
      const size = mini ? 80 : 160;
      const stroke = mini ? 8 : 16;
      const r = mini ? (size/2 - stroke/2 - 2) : (size/2 - stroke/2 - 4);
      const r2 = mini ? (r - 4) : (r - 8);
      
      if(isLeave) {
        return `<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}"><circle cx="${size/2}" cy="${size/2}" r="${r}" stroke="#e74c3c" stroke-width="${stroke}" fill="none"/><circle cx="${size/2}" cy="${size/2}" r="${r2}" fill="#ffe0e0"/></svg>`;
      } else {
        // 9~18시(9시간) = 270도(3/4원)
        const startX = size/2;
        const startY = size/2 - r;
        const endX = size - stroke/2;
        const endY = size/2;
        const largeArcFlag = 1;
        return `<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}"><circle cx="${size/2}" cy="${size/2}" r="${r}" stroke="#eee" stroke-width="${stroke}" fill="none"/><path d="M${startX} ${startY} A${r} ${r} 0 ${largeArcFlag} 1 ${endX} ${endY}" stroke="#218838" stroke-width="${stroke}" fill="none" stroke-linecap="round"/><circle cx="${size/2}" cy="${size/2}" r="${r2}" fill="#ffe0e0"/></svg>`;
      }
    }
  }
};
</script>
 