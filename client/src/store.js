import { compose, createStore, applyMiddleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import { createBrowserHistory } from 'history';
import { routerMiddleware } from 'connected-react-router';

import createRootReducer from './reducers';
import rootSaga from './sagas';

export const history = createBrowserHistory();
const sagaMiddleware = createSagaMiddleware();

export default preloadedState => {
  const store = createStore(
    createRootReducer(history),
    preloadedState,
    compose(applyMiddleware(routerMiddleware(history), sagaMiddleware))
  );
  sagaMiddleware.run(rootSaga);
  return store;
};
