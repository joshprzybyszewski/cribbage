import { combineReducers, createStore } from 'redux';
import { devToolsEnhancer } from 'redux-devtools-extension';

import { reducer as alerts } from '../app/containers/Alert/slice';
import { reducer as game } from '../app/containers/Game/slice';
import { reducer as home } from '../app/containers/Home/slice';
import { reducer as auth } from '../auth/slice';

/* Create root reducer, containing all features of the application */
const rootReducer = combineReducers({
    alerts,
    auth,
    home,
    game,
});

const store = createStore(
    rootReducer,
    /* preloadedState, */ devToolsEnhancer({}),
);

export type RootState = ReturnType<typeof rootReducer>;
export default store;
