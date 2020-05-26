import { auth } from '../sagas/types';

const actions = auth.reducer;

const initialState = {
  id: '',
  name: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case actions.LOGIN:
    case actions.REGISTER:
      return { ...state, ...payload };
    case actions.REGISTER_FAILED:
    case actions.LOGIN_FAILED:
    case actions.LOGOUT:
      return { id: '', name: '' };
    default:
      return state;
  }
};
