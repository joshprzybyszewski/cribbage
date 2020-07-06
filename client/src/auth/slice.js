import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentUser: {
    id: '',
    name: '',
  },
  activeGames: {},
  loading: false,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    login: {
      reducer: (state, action) => {
        state.loading = true;
        state.currentUser.id = action.payload.id;
      },
      prepare: (id, history) => {
        return { payload: { id, history } };
      },
    },
    loginSuccess(state, action) {
      state.loading = false;
      state.currentUser.id = action.payload.id;
      state.currentUser.name = action.payload.name;
    },
    loginFailed(state) {
      state.loading = false;
      state.currentUser = { id: '', name: '' };
    },
    register: {
      reducer: (state, action) => {
        state.loading = true;
        state.currentUser.id = action.payload.id;
        state.currentUser.name = action.payload.name;
      },
      prepare: (id, name, history) => {
        return { payload: { id, name, history } };
      },
    },
    registerSuccess(state, action) {
      state.loading = false;
      state.currentUser.id = action.payload.id;
      state.currentUser.name = action.payload.name;
    },
    registerFailed(state) {
      state.loading = false;
      state.currentUser = { id: '', name: '' };
    },
    logout: {
      reducer: state => {
        state.loading = false;
        state.currentUser = { id: '', name: '' };
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    refreshActiveGames(state, action) {
      if (state.currentUser.id !== action.payload.id) {
        // what should we do when refreshing with an ID we do not expect?
        throw `bad user id: expected "${state.currentUser.id}", got "${action.payload.id}"`;
      }
    },
    gotActiveGames(state, action) {
      if ( state.currentUser.id === action.payload.player.id ) {
        state.activeGames = action.payload.activeGames;
      }
    },
  },
});

export const { actions, reducer, name: sliceKey } = authSlice;
