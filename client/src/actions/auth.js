import { LOGIN_SUCCESS, REGISTER_SUCCESS } from './types';

export const login = (username) => async (dispatch) => {
  dispatch({
    type: LOGIN_SUCCESS,
    payload: username,
  });
};

export const register = (username) => async (dispatch) => {
  dispatch({
    type: REGISTER_SUCCESS,
    payload: username,
  });
};
