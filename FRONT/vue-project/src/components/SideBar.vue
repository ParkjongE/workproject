<template>
  <nav class="sidebar">
    <div class="user-info" v-if="currentUser">
      <div class="user-name">{{ currentUser.name }}</div>
      <div class="user-role">{{ currentUser.role }} | {{ currentUser.department }}</div>
      <button class="logout-btn" @click="logout">로그아웃</button>
    </div>
    <ul>
      <li :class="{ active: isActiveMenu('commute') }" @click="toggleSubmenu('commuteSubmenu')">
        출퇴근/근태 관리
        <ul class="submenu" :class="{ show: isSubmenuVisible('commuteSubmenu') }">
          <li :class="{ active: isActivePage('/commute') }" @click.stop="navigateTo('/commute')">출퇴근 기록</li>
          <li :class="{ active: isActivePage('/leave-apply') }" @click.stop="navigateTo('/leave-apply')">휴가/조퇴/외출 신청</li>
          <li :class="{ active: isActivePage('/my-attendance') }" @click.stop="navigateTo('/my-attendance')">내 근태내역</li>
        </ul>
      </li>
      <li :class="{ active: isActiveMenu('schedule') }" @click="toggleSubmenu('scheduleSubmenu')">
        내 일정관리
        <ul class="submenu" :class="{ show: isSubmenuVisible('scheduleSubmenu') }">
          <li :class="{ active: isActivePage('/my-schedule') }" @click.stop="navigateTo('/my-schedule')">나의 일정 수정</li>
          <li :class="{ active: isActivePage('/meeting-requests') }" @click.stop="navigateTo('/meeting-requests')">협업/미팅 요청 관리</li>
          <li :class="{ active: isActivePage('/work-records') }" @click.stop="navigateTo('/work-records')">내 업무 기록</li>
        </ul>
      </li>
      <li :class="{ active: isActiveMenu('share') }" @click="toggleSubmenu('shareSubmenu')">
        일정 공유
        <ul class="submenu" :class="{ show: isSubmenuVisible('shareSubmenu') }">
          <li :class="{ active: isActivePage('/team-schedule') }" @click.stop="navigateTo('/team-schedule')">팀 일정 공유</li>
          <li :class="{ active: isActivePage('/personal-schedule') }" @click.stop="navigateTo('/personal-schedule')">개인 일정 공유</li>
          <li :class="{ active: isActivePage('/interaction-graph') }" @click.stop="navigateTo('/interaction-graph')">상호작용 그래프</li>
          <li :class="{ active: isActivePage('/related-work') }" @click.stop="navigateTo('/related-work')">연관업무 기록 조회</li>
        </ul>
      </li>
      <li :class="{ active: isActiveMenu('org') }" @click="toggleSubmenu('orgSubmenu')">
        조직관리
        <ul class="submenu" :class="{ show: isSubmenuVisible('orgSubmenu') }">
          <li :class="{ active: isActivePage('/org-chart') }" @click.stop="navigateTo('/org-chart')">조직도</li>
          <li :class="{ active: isActivePage('/dep-work') }" @click.stop="navigateTo('/dep-work')">부서별 업무 현황 조회</li>
          <li :class="{ active: isActivePage('/org-member') }" @click.stop="navigateTo('/org-member')">구성원 관리</li>
        </ul>
      </li>
    </ul>
  </nav>
</template>

<script>
import api from '@/plugins/axios'

export default {
  name: 'SideBar',
  data() {
    return {
      visibleSubmenus: {
        commuteSubmenu: true,
        scheduleSubmenu: false,
        shareSubmenu: false,
        orgSubmenu: false
      },
      currentUser: null
    };
  },
  mounted() {
    // 현재 경로에 따라 적절한 서브메뉴를 자동으로 열어줍니다
    const path = this.$route.path;
    if (path.includes('/commute') || path.includes('/leave-apply') || path.includes('/my-attendance')) {
      this.visibleSubmenus.commuteSubmenu = true;
    } else if (path.includes('/my-schedule') || path.includes('/meeting-requests') || path.includes('/work-records')) {
      this.visibleSubmenus.scheduleSubmenu = true;
    } else if (path.includes('/team-schedule') || path.includes('/personal-schedule') || 
               path.includes('/interaction-graph') || path.includes('/related-work') ||
               path.includes('/schedule-sharing')) {
      this.visibleSubmenus.shareSubmenu = true;
    } else if (path.includes('/org-chart') || path.includes('/dep-work') || path.includes('/org-member')) {
      this.visibleSubmenus.orgSubmenu = true;
    }
    
    // 로그인한 사용자 정보 불러오기
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        this.currentUser = JSON.parse(userStr);
        console.log('사용자 정보 로드됨:', this.currentUser);
      } catch (e) {
        console.error('사용자 정보 파싱 오류:', e);
      }
    }
  },
  methods: {
    toggleSubmenu(submenu) {
      this.visibleSubmenus[submenu] = !this.visibleSubmenus[submenu];
    },
    isSubmenuVisible(submenu) {
      return this.visibleSubmenus[submenu];
    },
    navigateTo(route) {
      this.$router.push(route);
    },
    isActivePage(route) {
      return this.$route.path === route;
    },
    isActiveMenu(menu) {
      const path = this.$route.path;
      if (menu === 'commute') {
        return path.includes('/commute') || path.includes('/leave-apply') || path.includes('/my-attendance');
      } else if (menu === 'schedule') {
        return path.includes('/my-schedule') || path.includes('/meeting-requests') || path.includes('/work-records');
      } else if (menu === 'share') {
        return path.includes('/team-schedule') || path.includes('/personal-schedule') || 
               path.includes('/interaction-graph') || path.includes('/related-work') ||
               path.includes('/schedule-sharing');
      } else if (menu === 'org') {
        return path.includes('/org-chart') || path.includes('/dep-work') || path.includes('/org-member');
      }
      return false;
    },
    async logout() {
      try {
        // 서버에 로그아웃 요청
        await api.post('/api/logout');
      } catch (error) {
        console.error('로그아웃 요청 실패:', error);
      } finally {
        // 로컬 스토리지 정리
        localStorage.removeItem('user');
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('rememberMe');
        // 로그인 페이지로 리다이렉트
        this.$router.push('/');
      }
    }
  }
};
</script>

<style scoped>
.sidebar {
  width: 220px;
  background: #f5f5f5;
  border-right: 1px solid #e0e0e0;
  min-height: 700px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  padding: 0;
}
.user-info {
  padding: 20px;
  background-color: #e8f5e9;
  border-bottom: 1px solid #c8e6c9;
  margin-bottom: 10px;
}
.user-name {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 5px;
}
.user-role {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}
.logout-btn {
  background-color: #f44336;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}
.logout-btn:hover {
  background-color: #d32f2f;
}
.sidebar ul {
  list-style: none;
  padding: 36px 0 0 20px;
  margin: 0;
}
.sidebar li {
  margin: 18px 0;
  color: #333;
  cursor: pointer;
  font-size: 17px;
  font-weight: normal;
  transition: font-weight 0.2s;
  user-select: none;
}
.sidebar li.active {
  font-weight: bold;
  color: #218838;
}
.sidebar li:first-child {
  font-size: 17px;
  font-weight: normal;
  color: #333;
  margin-bottom: 18px;
}
.sidebar .submenu {
  display: none;
  margin-left: 12px;
  padding-top: 0;
}
.sidebar .submenu.show {
  display: block;
}
.sidebar .submenu li {
  font-size: 16px;
  color: #444;
  margin: 8px 0;
  font-weight: normal;
}
.sidebar .submenu li.active {
  font-weight: bold;
  color: #218838;
}
@media (max-width: 1200px) {
  .sidebar {
    width: 100%;
    min-height: unset;
  }
}
</style> 