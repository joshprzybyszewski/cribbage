import { REMOVE_ALERT, SET_ALERT_WORKER } from '../sagas/types';

const initialState = [];

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case SET_ALERT_WORKER:
      return [...state, payload];
    case REMOVE_ALERT:
      return state.filter(a => a.id !== payload);
    default:
      return state;
  }
};
