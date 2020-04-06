import { LOGIN_SUCCESS, REGISTER_SUCCESS } from './types';

export const login = (username) => (dispatch) => {
  alert(`LOGIN: ${username}`);
  dispatch({
    type: LOGIN_SUCCESS,
    payload: username,
  });
};

export const register = (username) => (dispatch) => {
  alert(`REGISTER: ${username}`);
  dispatch({
    type: REGISTER_SUCCESS,
    payload: username,
  });
};
