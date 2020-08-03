import { createSelector } from '@reduxjs/toolkit';
import { initialState } from 'app/containers/Game/slice';

const selectDomain = state => state.game || initialState;

export const selectCurrentGameID = createSelector(
  [selectDomain],
  gameState => gameState.currentGameID,
);

export const selectCurrentGame = createSelector(
  [selectDomain],
  gameState => gameState.currentGame,
);

export const selectIsLoading = createSelector(
  [selectDomain],
  gameState => gameState.isLoading,
);

export const selectSelectedCards = createSelector(
  [selectDomain],
  gameState => gameState.selectedCards,
);
