import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import { useStore } from './store.ts'
import PageCreate from '../pages/PageCreate.vue'
import PageHome from '../pages/PageHome.vue'
import PageList from '../pages/PageList.vue'
import PageLogin from '../pages/PageLogin.vue'
import PageNotFound from '../pages/PageNotFound.vue'
import PageShare from '../pages/PageShare.vue'

declare module 'vue-router' {
  interface RouteMeta {
    auth?: boolean
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: PageHome
  },
  {
    path: '/create',
    component: PageCreate,
    meta: {
      auth: true
    }
  },
  {
    path: '/login',
    component: PageLogin
  },
  {
    path: '/shares',
    component: PageList,
    meta: {
      auth: true
    }
  },
  {
    path: '/shares/:name',
    component: PageShare
  },
  {
    path: '/:pathMatch(.*)*',
    component: PageNotFound
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(async (to) => {
  const store = useStore()
  if (!store.userinfoLoaded) {
    if (to.meta?.auth !== undefined) {
      await store.fetchUserinfo()
    } else {
      store.fetchUserinfo() // auth is not required, does not await
    }
  }
  if (to.path !== '/' && to.meta?.auth && !store.loggedIn) {
    return '/'
  }
  return true
})

export default router
