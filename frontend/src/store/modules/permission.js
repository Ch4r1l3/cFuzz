import { adminRoutes, constantRoutes, lastRoutes } from '@/router'

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes).concat(lastRoutes)
  }
}

const actions = {
  generateRoutes({ commit }, isAdmin) {
    return new Promise(resolve => {
      let accessedRoutes = []
      if (isAdmin) {
        accessedRoutes = adminRoutes
      }
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
