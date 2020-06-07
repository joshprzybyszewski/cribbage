import { game } from '../sagas/types';

const actions = game.reducer;

const initialState = {
  gameID: '',
};

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case actions.VIEW_GAME:
      return { ...state, ...payload };
    case actions.VIEW_GAME_FAILED:
    case actions.EXIT_GAME:
      return { gameID: '' };
    default:
      return state;
  }
};
