import { auth, alert, game } from './types';

const addAlert = (msg, type) => ({
  type: alert.ADD_ALERT,
  payload: { msg, type },
});

export const alertActions = {
  addAlert,
};

const logout = () => ({
  type: auth.LOGOUT,
});

const login = id => ({
  type: auth.LOGIN,
  payload: id,
});

const register = (id, name) => ({
  type: auth.REGISTER,
  payload: { id, name },
});

export const authActions = {
  logout,
  login,
  register,
};

const viewGame = g => ({
  type: game.VIEW_GAME,
  payload: g,
});

export const gameActions = {
  viewGame,
}