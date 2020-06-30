import { combineReducers } from '@reduxjs/toolkit';
import { connectRouter } from 'connected-react-router';
import { createBrowserHistory } from 'history';

export const history = createBrowserHistory();

export const createReducer = injectedReducers => {
  if (Object.keys(injectedReducers).length === 0) {
    return state => state;
  }
  return combineReducers({
    router: connectRouter(history),
    ...injectedReducers,
  });
};
