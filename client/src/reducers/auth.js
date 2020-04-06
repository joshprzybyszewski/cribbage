import { LOGIN_SUCCESS, REGISTER_SUCCESS } from './types';

const initialState = {
  user: '',
  loading: true,
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN_SUCCESS:
      return { ...state, user: payload, loading: false };
    case REGISTER_SUCCESS:
      return { ...state, user: payload, loading: false };
    default:
      return state;
  }
};
