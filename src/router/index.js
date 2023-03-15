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
      {
        path: '/login',
        name: 'Login',
        beforeEnter() {
          let url = '/sn/login'
          if (process.env.NODE_ENV == "development") {
            url = 'https://plateau.fake-slothninja.com:8091' + url
          }
          window.location.replace(url)
        }
      },
      {
        path: '/logout',
        name: 'Logout',
        beforeEnter() {
          let url = '/sn/logout'
          if (process.env.NODE_ENV == "development") {
            url = 'https://plateau.fake-slothninja.com:8091' + url
          }
          window.location.replace(url)
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
