const state = {
  count: 0,
  currentPage: 1,
  pageSize: 25
}

const mutations = {
  SET_COUNT: (state, count) => {
    state.count = count
  },
  SET_CURRENT_PAGE: (state, currentPage) => {
    state.currentPage = currentPage
  }
}

const actions = {
  setCount({ commit }, count) {
    return new Promise((resolve, reject) => {
      if (count < 0) {
        reject('count cannot small than 0')
      } else {
        commit('SET_COUNT', count)
        resolve()
      }
    })
  },
  setCurrentPage({ commit, state }, currentPage) {
    return new Promise((resolve, reject) => {
      if (Math.trunc(state.count / state.pageSize) + 1 < currentPage) {
        reject('currentPage error')
      } else {
        commit('SET_CURRENT_PAGE', currentPage)
        resolve()
      }
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
