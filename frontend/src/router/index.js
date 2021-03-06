import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * all roles can be accessed
 */
const docUrl = location.protocol.concat('//').concat(location.host).concat('/api/docs/swagger.json')

export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/task'
  },
  {
    path: '/task',
    component: Layout,
    redirect: '/task/list',
    children: [
      {
        path: 'list',
        name: 'listTask',
        component: () => import('@/views/task/list'),
        meta: { title: 'Task', icon: 'table' }
      },
      {
        path: 'create',
        name: 'createTask',
        component: () => import('@/views/task/create'),
        meta: { title: 'Create' },
        hidden: true
      },
      {
        path: 'edit/:id(\\d+)',
        name: 'editTask',
        component: () => import('@/views/task/edit'),
        meta: { title: 'Edit' },
        hidden: true
      },
      {
        path: 'detail/:id(\\d+)',
        name: 'taskDetail',
        component: () => import('@/views/task/detail'),
        meta: { title: 'Detail' },
        hidden: true
      }

    ]
  },

  {
    path: '/image',
    component: Layout,
    redirect: '/image/list',
    children: [
      {
        path: 'list',
        name: 'listImage',
        component: () => import('@/views/image/list'),
        meta: { title: 'Image', icon: 'edit' }
      },
      {
        path: 'create',
        name: 'createImage',
        component: () => import('@/views/image/create'),
        meta: { title: 'Create' },
        hidden: true
      },
      {
        path: 'edit/:id(\\d+)',
        name: 'editImage',
        component: () => import('@/views/image/edit'),
        meta: { title: 'Edit' },
        hidden: true
      }

    ]
  },
  {
    path: '/storageItem',
    component: Layout,
    redirect: '/storageItem/list',
    children: [
      {
        path: 'list',
        name: 'listStorageItem',
        component: () => import('@/views/storageItem/list'),
        meta: { title: 'StorageItem', icon: 'documentation' }
      },
      {
        path: 'create',
        name: 'createStorageItem',
        component: () => import('@/views/storageItem/create'),
        meta: { title: 'Create' },
        hidden: true
      },
      {
        path: 'edit/:id(\\d+)',
        name: 'editStorageItem',
        component: () => import('@/views/storageItem/edit'),
        meta: { title: 'Edit' },
        hidden: true
      }

    ]
  }

]

export const adminRoutes = [
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    children: [
      {
        path: 'list',
        name: 'listUser',
        component: () => import('@/views/user/list'),
        meta: { title: 'User', icon: 'user' }
      },
      {
        path: 'create',
        name: 'createUser',
        component: () => import('@/views/user/create'),
        meta: { title: 'Create' },
        hidden: true
      }
    ]
  }
]

export const lastRoutes = [
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: encodeURI('http://petstore.swagger.io/?url=' + docUrl),
        meta: { title: 'Api Document', icon: 'link' }
      }
    ]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
