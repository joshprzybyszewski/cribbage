import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentUser: {
    id: '',
    name: '',
  },
  loading: true,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    login(state, action) {
      state = initialState;
      state.currentUser.id = action.payload;
    },
    loginSuccess(state, action) {
      state.loading = false;
      state.currentUser.id = action.payload.id;
      state.currentUser.name = action.payload.name;
    },
    loginFailed(state) {
      state = initialState;
      state.loading = false;
    },
    register(state, action) {
      state = initialState;
      state.currentUser.id = action.payload.id;
      state.currentUser.name = action.payload.name;
    },
    registerSuccess(state, action) {
      state.loading = false;
      state.currentUser.id = action.payload.id;
      state.currentUser.name = action.payload.name;
    },
    registerFailed(state) {
      state = initialState;
      state.loading = false;
    },
    logout(state) {
      state = initialState;
      state.loading = false;
    },
  },
});

export const { actions, reducer, name: sliceKey } = authSlice;
