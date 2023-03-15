// Composables
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "home" */ '@/views/Home.vue'),
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    beforeEnter() {
      window.location.replace('https://plateau.fake-slothninja.com:8091/sn/login')
    }
  },
  {
    path: '/logout',
    name: 'Logout',
    beforeEnter() {
      window.location.replace('https://plateau.fake-slothninja.com:8091/sn/logout')
    }
  },

]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
