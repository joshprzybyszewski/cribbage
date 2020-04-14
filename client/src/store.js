import { createStore, applyMiddleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import rootReducer from './reducers';
import rootSaga from './sagas';

const sagaMiddleware = createSagaMiddleware();

const initialState = {};
const middleware = [sagaMiddleware];
const store = createStore(
  rootReducer,
  initialState,
  applyMiddleware(...middleware)
);
sagaMiddleware.run(rootSaga);

export default store;
