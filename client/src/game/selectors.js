import { createSelector } from '@reduxjs/toolkit';
import { initialState } from './slice';

const selectDomain = state => state.game || initialState;

export const selectCurrentGameID = createSelector(
  [selectDomain],
  gameState => gameState.currentGameID,
);

export const selectCurrentGame = createSelector(
    [selectDomain],
    gameState => gameState.currentGame,
  );