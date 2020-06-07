import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import auth from './auth';
import alert from './alert';
import game from './game';

const createRootReducer = history =>
  combineReducers({
    router: connectRouter(history),
    auth,
    alert,
    game,
  });

export default createRootReducer;
