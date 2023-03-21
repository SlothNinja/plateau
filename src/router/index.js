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
      }
    ]
  },
  {
    path: '/new',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'NewInvitation',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "new" */ '@/views/NewInvitation.vue'),
      }
    ]
  },
  {
    path: '/join',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'InvitationIndex',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "new" */ '@/views/InvitationIndex.vue'),
      }
    ]
  },
  {
    path: '/games/:status',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'GameIndex',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "new" */ '@/views/GameIndex.vue'),
      }
    ]
  },
  {
    path: '/game/:id',
    component: () => import('@/layouts/game/Default.vue'),
    children: [
      {
        path: '',
        name: 'GameShow',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "new" */ '@/views/Game.vue'),
      }
    ]
  },
  {
    path: '/user/:id',
    name: 'User',
    beforeEnter(to) {
      let url = `/#/show/${to.params.id}`
      if (process.env.NODE_ENV == "development") {
        url = 'https://user.fake-slothninja.com:8087' + url
      } else {
        url = 'https://user.slothninja.com' + url
      }
      window.location.replace(url)
    }
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
  },
  {
    path: '/webpublished',
    name: 'WebPublished',
    beforeEnter() {
      let url = 'https://boardgamegeek.com/boardgame/339349/le-plateau'
      window.location.replace(url)
    }
  },
  {
    path: '/rules',
    name: 'Rules',
    beforeEnter() {
      let url = 'https://boardgamegeek.com/filepage/223027/le-plateau-rules'
      window.location.replace(url)
    }
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
