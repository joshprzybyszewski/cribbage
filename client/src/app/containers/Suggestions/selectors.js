import { createSelector } from '@reduxjs/toolkit';

const selectDomain = state => state.suggestions;

export const selectSuggestions = createSelector(
    [selectDomain],
    s => s.suggestedHands,
);

export const selectHandCards = createSelector(
    [selectDomain],
    s => s.handCards,
);