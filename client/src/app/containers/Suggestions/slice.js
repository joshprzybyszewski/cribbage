import { nanoid } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
    handCards: [
        'AH', 'KH', '5H', 'JH', '8C', '3D',
    ],
    suggestedHands: [{
        hand: ['AH', 'KH', '5H', 'JH', ],
        throw: ['8C', '3D', ],
        handPts: { avg: 5, median: 10 },
        cribPts: { avg: 2, median: 8 },
    }, {
        hand: ['AH', 'KH', '8C', '3D', ],
        throw: ['5H', 'JH', ],
        handPts: { avg: 7, median: 15 },
        cribPts: { avg: 4, median: 6 },
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