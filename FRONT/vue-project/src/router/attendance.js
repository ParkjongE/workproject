const attendanceRoutes = [
  {
    path: '/commute',
    name: 'commute',
    component: () => import('../views/attendance/Commute.vue')
  },
  {
    path: '/my-attendance',
    name: 'myAttendance',
    component: () => import('../views/attendance/MyAttendance.vue')
  },
  {
    path: '/leave-apply',
    name: 'leaveApply',
    component: () => import('../views/attendance/LeaveApply.vue')
  }
]

export default attendanceRoutes 