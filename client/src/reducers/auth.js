import { LOGIN, REGISTER } from '../sagas/types';

const initialState = {
  user: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`LOGIN: ${payload}`);
      return { ...state, user: payload };
    case REGISTER:
      console.log(`REGISTER: ${payload}`);
      return { ...state, user: payload };
    default:
      return state;
  }
};
