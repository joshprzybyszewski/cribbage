import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { Card, Game } from './models';

export interface DealAction {
    numShuffles: number;
}

export interface SelectCardsAction {
    selectedCards: Card[];
}

export type CribAction = SelectCardsAction;
export type PegAction = SelectCardsAction;

export interface CutAction {
    percCut: number;
}

export interface CountAction {
    points: number;
}
export type CountHandAction = CountAction;
export type CountCribAction = CountAction;

export type GameAction =
    | DealAction
    | CribAction
    | PegAction
    | CutAction
    | CountHandAction
    | CountCribAction;

export interface GameState {
    currentGameID: number;
    currentGame: Game;
    selectedCards: Card[];
    currentAction: GameAction;
    loading: boolean;
}

export const initialState: GameState = {
    currentGameID: 0,
    currentGame: {
        id: 0,
        teams: [],
        phase: 'unknownPhase',
        current_peg: 0,
        blocking_players: {},
        current_dealer: '',
        cut_card: {
            name: 'AC',
            value: 1,
            suit: 'Clubs',
        },
        hands: {},
        pegged_cards: [],
        crib: [],
    },
    selectedCards: [],
    currentAction: {
        numShuffles: 0,
        selectedCards: [],
        percCut: 0.5,
        points: -1,
    },
    loading: true,
};

const gameSlice = createSlice({
    name: 'game',
    initialState,
    reducers: {
        setLoading(state, action: PayloadAction<boolean>) {
            return {
                ...state,
                loading: action.payload,
            };
        },
        setGameID(state, action: PayloadAction<number>) {
            return {
                ...state,
                currentGameID: action.payload,
            };
        },
        setGame(state, action: PayloadAction<Game>) {
            return {
                ...state,
                currentGame: action.payload,
                currentGameID: action.payload.id,
            };
        },
        exitGame() {
            return initialState;
        },
        toggleSelectedCard(state, action: PayloadAction<Card>) {
            const card = action.payload;
            const cardsAreEqual = (c1: Card, c2: Card) =>
                c1.suit === c2.suit && c1.value === c2.value;
            if (state.selectedCards.some(c => cardsAreEqual(c, card))) {
                return {
                    ...state,
                    selectedCards: state.selectedCards.filter(
                        c => !cardsAreEqual(c, card),
                    ),
                };
            }
            return {
                ...state,
                selectedCards: [...state.selectedCards, card],
            };
        },
        clearSelectedCards(state) {
            return {
                ...state,
                selectedCards: [],
            };
        },
    },
});

export const { actions, reducer, name: sliceKey } = gameSlice;
