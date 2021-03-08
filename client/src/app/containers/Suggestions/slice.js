import { nanoid } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
    handCards: [
        'AH', 'KH', '5H', 'JH', '8C', '3D',
    ],
    loading: false,
    suggestedHands: [{
        hand: ['AH', 'KH', '5H', 'JH', ],
        toss: ['8C', '3D', ],
        handPts: { avg: 5, median: 10 },
        cribPts: { avg: 2, median: 8 },
    }, {
        hand: ['AH', 'KH', '8C', '3D', ],
        toss: ['5H', 'JH', ],
        handPts: { avg: 7, median: 15 },
        cribPts: { avg: 4, median: 6 },
    }]
};

const suggestionsSlice = createSlice({
    name: 'suggestions',
    initialState,
    reducers: {
        getHandSuggestion: {
            reducer: (state, action) => {
                state.loading = true;
            },
        },
        setSuggestionResult: {
            reducer: (state, action) => {
                state.loading = false;
                state.suggestedHands = action.payload.data;
            },
        },
        updateCard: {
            reducer: (state, action) => {
                const selCard = action.payload.card;
                if (!state.handCards.includes(selCard)) {
                    throw `hand doesn't include card: ${selCard}`;
                }
                if (!action.payload.newCard) {
                    throw `needs newCard`;
                }
                let cpy = [...state.handCards];
                cpy.splice(
                    cpy.indexOf(selCard),
                    1,
                    action.payload.newCard);
                state.handCards = cpy;
            },
        },
    },
});

export const { actions, reducer, name: sliceKey } = suggestionsSlice;