import { alert } from '../sagas/types';

const actions = alert.reducer;

const initialState = [];

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case actions.ADD_ALERT:
      return [payload, ...state];
    case actions.REMOVE_ALERT:
      return state.filter(a => a.id !== payload);
    default:
      return state;
  }
};
