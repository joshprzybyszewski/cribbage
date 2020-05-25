import { LOGIN, REGISTER, REGISTER_FAILED, LOGIN_FAILED } from '../sagas/types';

const initialState = {
  id: '',
  name: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`LOGIN: ${payload.id}`);
      return { ...state, ...payload };
    case REGISTER:
      console.log(`REGISTER: ${payload.id} with displayName = ${payload.name}`);
      return { ...state, ...payload };
    case REGISTER_FAILED:
    case LOGIN_FAILED:
      console.log('AUTH FAILED WITH ERROR');
      console.log(payload);
      return { id: '', name: '' };
    default:
      return state;
  }
};
