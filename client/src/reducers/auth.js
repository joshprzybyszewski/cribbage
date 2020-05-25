import { LOGIN, REGISTER, REGISTER_FAILED } from '../sagas/types';

const initialState = {
  id: '',
  name: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case LOGIN:
      console.log(`LOGIN: ${payload.username}`);
      return { ...state, id: payload.username };
    case REGISTER:
      console.log(`REGISTER: ${payload.id} with displayName = ${payload.name}`);
      return { ...state, ...payload };
    case REGISTER_FAILED:
      console.log('REGISTER FAILED WITH ERROR');
      console.log(payload);
      return { id: '', name: '' };
    default:
      return state;
  }
};
