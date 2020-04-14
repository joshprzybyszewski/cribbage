import { LOGIN } from '../sagas/types';

const initialState = {
  user: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`PAYLOAD: ${payload}`);
      return { ...state, user: payload };
    default:
      return state;
  }
};
