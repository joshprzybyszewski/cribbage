import { LOGIN, REGISTER } from '../sagas/types';

const initialState = {
  user: '',
  displayName: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`LOGIN: ${payload.username}`);
      return { ...state, user: payload.username };
    case REGISTER:
      console.log(
        `REGISTER: ${payload.username} with displayName = ${payload.displayName}`
      );
      return { ...state, ...payload };
    default:
      return state;
  }
};
