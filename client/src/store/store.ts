import { combineReducers, createStore } from 'redux';
import { devToolsEnhancer } from 'redux-devtools-extension';

import { reducer as alerts } from '../app/containers/Alert/slice';

/* Create root reducer, containing all features of the application */
const rootReducer = combineReducers({
    alerts,
});

const store = createStore(
    rootReducer,
    /* preloadedState, */ devToolsEnhancer({}),
);

export default store;
