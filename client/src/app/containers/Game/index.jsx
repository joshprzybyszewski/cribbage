import React from 'react';
import { useSelector } from 'react-redux';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser, selectActiveGames } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';

const Game = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGames = useSelector(selectActiveGames);
  
  // event handlers
  const onRefreshActiveGames = () => {
      // TODO change this to refreshing the current game
    dispatch(actions.refreshActiveGames(currentUser.id));
  };

  return <div>
      This will be a page for the game of a user.
      <button
          onClick={onRefreshActiveGames}
          className='hover:text-white'
        >
          Refresh
        </button>
    </div>;
};

export default Game;
