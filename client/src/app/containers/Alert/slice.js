import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  alerts: [],
};

const alertSlice = createSlice({
  name: 'alert',
  initialState,
  reducers: {
    requestAlert(state) {
      // idk if we need to do this
      state = state;
    },
    addAlert(state, action) {
      state.alerts = [...state.alerts, action.payload];
    },
    removeAlert(state, action) {
      state.alerts.filter(a => a.id !== action.payload);
    },
  },
});

export const { actions, reducer, name: sliceKey } = alertSlice;
