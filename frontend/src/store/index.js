import Vue from 'vue'
import Vuex from 'vuex'
import getters from './getters'
import app from './modules/app'
import settings from './modules/settings'
import deployment from './modules/deployment'
import storageItem from './modules/storageItem'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    app,
    settings,
    deployment,
    storageItem
  },
  getters
})

export default store
