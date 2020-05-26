import {
  LOGIN,
  LOGOUT,
  REGISTER,
  REGISTER_FAILED,
  LOGIN_FAILED,
} from '../sagas/types';

const initialState = {
  id: '',
  name: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      return { ...state, ...payload };
    case REGISTER:
      return { ...state, ...payload };
    case REGISTER_FAILED:
    case LOGIN_FAILED:
    case LOGOUT:
      return { id: '', name: '' };
    default:
      return state;
  }
};
