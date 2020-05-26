import { combineReducers } from 'redux';
import { connectRouter } from 'connected-react-router';

import auth from './auth';
import alert from './alert';

const createRootReducer = history =>
  combineReducers({
    router: connectRouter(history),
    auth,
    alert,
  });

export default createRootReducer;
