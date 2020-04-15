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
      console.log(`REGISTER ERRORED: ${payload}`);
      return state;
    case REGISTER_SUCCESS:
      console.log(`REGISTERED PLAYER: ${payload}`);
      return { ...state, ...payload };
    default:
      return state;
  }
};
