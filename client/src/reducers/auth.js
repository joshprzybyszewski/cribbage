import { auth } from '../sagas/types';

const actions = auth.reducer;

const initialState = {
  id: '',
  name: '',
  allgames: {},
  loggedIn: false,
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case actions.LOGIN:
    case actions.REGISTER:
      return { ...state, ...payload, loggedIn: true };
    case actions.REGISTER_FAILED:
    case actions.LOGIN_FAILED:
    case actions.LOGOUT:
      return { id: '', name: '', allgames: {}, loggedIn: false };
    default:
      return state;
  }
};
