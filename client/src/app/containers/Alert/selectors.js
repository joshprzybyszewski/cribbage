import { createSelector } from '@reduxjs/toolkit';

const selectDomain = state => state.alert;

export const selectAlerts = createSelector(
  [selectDomain],
  alert => alert.alerts,
);
