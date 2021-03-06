import { combineReducers, createStore } from 'redux';
import { devToolsEnhancer } from 'redux-devtools-extension';

import { reducer as alerts } from '../app/containers/Alert/slice';
import { reducer as auth } from '../auth/slice';

/* Create root reducer, containing all features of the application */
const rootReducer = combineReducers({
    alerts,
    auth,
});

const store = createStore(
    rootReducer,
    /* preloadedState, */ devToolsEnhancer({}),
);

export type RootState = ReturnType<typeof rootReducer>;
export default store;
