import {
  LOGIN_SUCCESS,
  LOGIN_FAIL,
  REGISTER_FAIL,
  REGISTER_SUCCESS,
} from '../sagas/types';

const initialState = {
  player: {},
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN_FAIL:
      return { ...state, player: {} };
    case LOGIN_SUCCESS:
      return { ...state, ...payload };
    case REGISTER_FAIL:
      return { ...state, player: {} };
    case REGISTER_SUCCESS:
      return { ...state, ...payload };
    default:
      return state;
  }
};
