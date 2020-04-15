import { LOGIN, REGISTER_FAIL, REGISTER_SUCCESS } from '../sagas/types';

const initialState = {
  player: {},
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`LOGIN: ${payload.username}`);
      return { ...state, user: payload.username };
    case REGISTER_FAIL:
      return { ...state, player: {} };
    case REGISTER_SUCCESS:
      return { ...state, ...payload };
    default:
      return state;
  }
};
