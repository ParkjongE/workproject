const scheduleSharingRoutes = [
  {
    path: '/team-schedule',
    name: 'teamSchedule',
    component: () => import('../views/schedulesharing/TeamSchedule.vue')
  },
  {
    path: '/personal-schedule',
    name: 'personalSchedule',
    component: () => import('../views/schedulesharing/ScheduleSharing.vue')
  },
  {
    path: '/interaction-graph',
    name: 'interactionGraph',
    component: () => import('../views/schedulesharing/InteractionGraph.vue')
  },
  {
    path: '/related-work',
    name: 'relatedWork',
    component: () => import('../views/schedulesharing/RelatedWork.vue')
  },
  // 이전 경로를 리다이렉트
  {
    path: '/schedule-sharing',
    redirect: '/personal-schedule'
  }
]

export default scheduleSharingRoutes 