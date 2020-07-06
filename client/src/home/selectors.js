import { createSelector } from '@reduxjs/toolkit';
import { initialState } from './slice';

const selectDomain = state => state.home || initialState;

export const selectActiveGames = (playerID) => createSelector(
  [selectDomain],
  homeState => homeState.activeGamesPlayerID === playerID ? homeState.activeGames : {},
);