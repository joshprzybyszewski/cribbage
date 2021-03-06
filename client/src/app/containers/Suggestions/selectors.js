import { createSelector } from '@reduxjs/toolkit';

const selectDomain = state => state.suggestions;

export const selectSuggestions = createSelector(
    [selectDomain],
    s => s.suggestions,
);