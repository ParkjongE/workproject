const myScheduleRoutes = [
  {
    path: '/my-schedule',
    name: 'Schedule',
    component: () => import('../views/myschedule/Schedule/Schedule.vue')
  },
  {
    path: '/meeting-requests',
    name: 'meetingRequests',
    component: () => import('../views/myschedule/MeetingRequests/MeetingRequests.vue')
  },
  {
    path: '/work-records',
    name: 'workRecords',
    component: () => import('../views/myschedule/WorkRecords/WorkRecords.vue')
  }
]

export default myScheduleRoutes 