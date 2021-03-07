import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { User } from '../../../auth/slice';

interface Player extends User {
    color: string;
}

/*
{
    "player": {
        "id": "123",
        "name": "abc"
    },
    "activeGames": [
        {
            "gameID": 3735260247,
            "players": [
                {
                    "id": "123",
                    "name": "abc",
                    "color": "blue"
                },
                {
                    "id": "124",
                    "name": "abc",
                    "color": "red"
                }
            ],
            "created": "0001-01-01T00:00:00Z",
            "lastMove": "0001-01-01T00:00:00Z"
        }
    ]
}
*/

export interface ActiveGameResponse {
    player: User;
    activeGames: ActiveGame[];
}

export interface ActiveGame {
    gameID: string;
    players: Player[];
    created: Date;
    lastMove: Date;
}

export const initialState: ActiveGameResponse = {
    player: {
        id: '',
        name: '',
    },
    activeGames: [],
};

const homeSlice = createSlice({
    name: 'home',
    initialState,
    reducers: {
        setActiveGamesPlayerID(state, action: PayloadAction<string>) {
            return {
                ...state,
                player: { id: action.payload, name: '' },
            };
        },
        setActiveGames(state, action: PayloadAction<ActiveGame[]>) {
            return {
                ...state,
                activeGames: action.payload,
            };
        },
    },
});

export const { actions, reducer, name: sliceKey } = homeSlice;
