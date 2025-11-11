import { createStore } from 'vuex'

export default {
  state: {
    user: {
      id: '',
      email: '',
      name: {
        first: '',
        last: ''
      },
      session: null,
      favorites: []
    },
    token: localStorage.getItem('token') || ''
  },
  mutations: {
    SET_USER(state, user) {
      state.user = user
      if (user.token) {
        state.token = user.token
        localStorage.setItem('token', user.token)
      }
    },
    CLEAR_USER(state) {
      state.user = {
        id: '',
        email: '',
        name: {
          first: '',
          last: ''
        },
        session: null,
        favorites: []
      }
      state.token = ''
      localStorage.removeItem('token')
    },
    ADD_FAVORITE(state, item) {
      state.user.favorites.push(item)
    },
    REMOVE_FAVORITE(state, index) {
      state.user.favorites.splice(index, 1)
    }
  },
  actions: {
    login({ commit }, userData) {
      commit('SET_USER', userData)
    },
    logout({ commit }) {
      commit('CLEAR_USER')
    },
    addFavorite({ commit }, item) {
      commit('ADD_FAVORITE', item)
    },
    removeFavorite({ commit }, index) {
      commit('REMOVE_FAVORITE', index)
    }
  },
  getters: {
    isAuthenticated: state => !!state.token || !!(state.user && state.user.id),
    userName: state => `${state.user.name.first} ${state.user.name.last}`.trim()
  }
}