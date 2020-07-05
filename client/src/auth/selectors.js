import { createSelector } from '@reduxjs/toolkit';
import { initialState } from './slice';

const selectDomain = state => state.auth || initialState;

export const selectCurrentUser = createSelector(
  [selectDomain],
  authState => authState.currentUser,
);

export const selectActiveGames = createSelector(
  [selectDomain],
  authState => authState.activeGames,
);

export const selectLoggedIn = createSelector(
  [selectDomain],
  authState => authState.currentUser.id !== '',
);

export const selectLoading = createSelector(
  [selectDomain],
  authState => authState.loading,
);
