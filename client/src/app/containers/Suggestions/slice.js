import { nanoid } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
    handCards: [
        'AH', 'KH', '5H', 'JH', '8C', '3D',
    ],
    suggestedHands: [{
        hand: ['AH', 'KH', '5H', 'JH', ],
        result: { avg: 5, median: 10 },
    }]
};

const suggestionsSlice = createSlice({
    name: 'suggestions',
    initialState,
    reducers: {
        setSuggestionResult: {
            reducer: (state, action) => {
                state.suggestedHands = action.payload.suggestions;
            },
        },
    },
});

export const { actions, reducer, name: sliceKey } = suggestionsSlice;