import { ADD_ALERT, REMOVE_ALERT } from '../sagas/types';

const initialState = [];

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case ADD_ALERT:
      return [payload, ...state];
    case REMOVE_ALERT:
      return state.filter(a => a.id !== payload);
    default:
      return state;
  }
};
