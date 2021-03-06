import { combineReducers } from '@reduxjs/toolkit';

export const createReducer = (injectedReducers = {}) => {
    if (Object.keys(injectedReducers).length === 0) {
        return state => state;
    }
    return combineReducers({
        ...injectedReducers,
    });
};
