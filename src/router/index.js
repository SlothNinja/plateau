// Composables
import { createRouter, createWebHistory } from 'vue-router'
import _includes from 'lodash/includes'

const sngGames = [ 'atf', 'confucius', 'indonesia', 'tammany', 'all' ]
const sngXGames = [ 'atf', 'confucius' ]

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
    path: '/sng-home',
    name: 'sng-home',
    beforeEnter(to) {
      const sngHome = import.meta.env.VITE_SNG_HOME
      window.location.replace(sngHome)
    }
  },
  {
    path: '/sng-games/:type/:status',
    name: 'sng-games',
    beforeEnter(to) {
      if (_includes(sngGames, to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        window.location.replace(`${sngHome}#/games/${to.params.status}/${to.params.type}`)
      } else {
        const gotHome = import.meta.env.VITE_GOT_HOME
        const tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        const plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/games/${to.params.status}`)
            break;
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/games/${to.params.status}`)
            break;
          case 'plateau':
            window.location.replace(`${plateauHome}games/${to.params.status}`)
        }
      }
    }
  },
  {
    path: '/sng-ugames/:uid/:status/:type',
    name: 'sng-ugames',
    beforeEnter(to) {
      const sngHome = import.meta.env.VITE_SNG_HOME
      window.location.replace(`${sngHome}#/ugames/${to.params.uid}/${to.params.status}/${to.params.type}`)
    }
  },
  {
    path: '/sng-new-game/:type',
    name: 'sng-new-game',
    beforeEnter(to) {
      if (_includes(sngGames, to.params.type)) {
        let sngHome = import.meta.env.VITE_SNG_HOME
        if (_includes(sngXGames, to.params.type)) {
          window.location.replace(`${sngHome}#/invitation/new/${to.params.type}`)
        } else {
          window.location.replace(`${sngHome}${to.params.type}/game/new`)
        }
      } else {
        let gotHome = import.meta.env.VITE_GOT_HOME
        let tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/invitation/new`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/invitation/new`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}new`)
        }
      }
    }
  },
  {
    path: '/sng-join-game/:type',
    name: 'sng-join-game',
    beforeEnter(to) {
      if (_includes(sngGames, to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        if (_includes(sngXGames, to.params.type)) {
          window.location.replace(`${sngHome}#/invitations/${to.params.type}`)
        } else {
          window.location.replace(`${sngHome}${to.params.type}/games/recruiting`)
        }
      } else {
        let gotHome = import.meta.env.VITE_GOT_HOME
        let tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/invitations`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/invitations`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}join`)
        }
      }
    }
  },
  {
    path: '/sng-ratings/:type',
    name: 'sng-ratings',
    beforeEnter(to) {
      if (_includes(sngGames, to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        window.location.replace(`${sngHome}ratings/show/${to.params.type}`)
      } else {
        let gotHome = import.meta.env.VUE_APP_GOT_HOME
        let tammanyHome = import.meta.env.VUE_APP_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/rank`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/rank`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}join`)
        }
      }
    }
  },
  {
    path: '/user/:id',
    name: 'User',
    beforeEnter(to) {
      const url = `${import.meta.env.VITE_USER_FRONTEND}user/${to.params.id}`
      window.location.replace(url)
    }
  },
  {
    path: '/login',
    name: 'Login',
    beforeEnter() {
      const url = `${import.meta.env.VITE_USER_BACKEND}sn/user/login`
      window.location.replace(url)
    }
  },
  {
    path: '/invitation/:action/:id',
    name: 'InvitationAction',
    beforeEnter(to) {
      const url = `${import.meta.env.VITE_PLATEAU_BACKEND}sn/invitation/${to.params.action}/${to.params.id}`
      window.location.replace(url)
    }
  },
  {
    path: '/logout',
    name: 'Logout',
    beforeEnter() {
      const query = `?redirect=${encodeURIComponent(import.meta.env.VITE_PLATEAU_FRONTEND)}`
      const url = `${import.meta.env.VITE_USER_BACKEND}sn/user/logout${query}`
      window.location.replace(url)
    }
  },
  {
    path: '/webpublished',
    name: 'WebPublished',
    beforeEnter() {
      const url = 'https://boardgamegeek.com/boardgame/339349/le-plateau'
      window.location.replace(url)
    }
  },
  {
    path: '/rules',
    name: 'Rules',
    beforeEnter() {
      const url = 'https://boardgamegeek.com/filepage/223027/le-plateau-rules'
      window.location.replace(url)
    }
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
