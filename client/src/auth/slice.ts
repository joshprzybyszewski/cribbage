import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface User {
    id: string;
    name: string;
}

interface AuthState {
    currentUser: User;
    loading: boolean;
}

export const initialState: AuthState = {
    currentUser: {
        id: '',
        name: '',
    },
    loading: false,
};

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setLoading(state, action: PayloadAction<boolean>) {
            return {
                ...state,
                loading: action.payload,
            };
        },
        setUser(state, action: PayloadAction<User>) {
            return {
                ...state,
                currentUser: action.payload,
            };
        },
        clearUser() {
            return initialState;
        },
    },
});

export const { actions, reducer, name: sliceKey } = authSlice;
