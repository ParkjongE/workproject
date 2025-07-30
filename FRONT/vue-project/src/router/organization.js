const organizationRoutes = [
  {
    path: '/org-chart',
    name: 'orgChart',
    component: () => import('../views/organization/OrgChart.vue')
  },
  {
    path: '/org-member',
    name: 'orgMember',
    component: () => import('../views/organization/OrgMember.vue')
  },
  {
    path: '/dep-work',
    name: 'depWork',
    component: () => import('../views/organization/DepWork.vue')
  }
]

export default organizationRoutes 