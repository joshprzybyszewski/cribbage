import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Card } from '../Game/models';

export interface CardUpdate {
    prev: Card;
    cur: Card;
}
export interface Stats {
    min: number;
    avg: number;
    median: number;
    max: number;
}

export interface TossSuggestion {
    hand: Card[];
    toss: Card[];
    handPts: Stats;
    cribPts: Stats;
}

export interface TossSuggestionState {
    handCards: Card[];
    suggestedHands: TossSuggestion[];
    loading: boolean;
}

export const initialState = {
    handCards: [{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    },{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    },{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    },{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    },{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    },{
        name: 'AC',
        value: 1,
        suit: 'Clubs',
    }],
    loading: false,
    suggestedHands: []
};

const suggestionsSlice = createSlice({
    name: 'suggestions',
    initialState,
    reducers: {
        setLoading(state, action: PayloadAction<boolean>) {
            return {
                ...state,
                loading: action.payload,
            };
        },
        setSuggestionResult(state, action: PayloadAction<TossSuggestion[]>) {
            return {
                ...state,
                loading: false,
                suggestedHands: action.payload,
            };
        },
        updateCard(state, action: PayloadAction<CardUpdate>) {
            const update = action.payload;
            const cardsAreEqual = (c1: Card, c2: Card) =>
                c1.suit === c2.suit && c1.value === c2.value;
            return {
                ...state,
                handCards: state.handCards.map(
                    (c) => cardsAreEqual(c, update.prev) ? update.cur : c
                ),
            };
        },
    },
});

export const { actions, reducer, name: sliceKey } = suggestionsSlice;