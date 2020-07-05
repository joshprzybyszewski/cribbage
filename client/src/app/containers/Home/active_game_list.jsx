import React from 'react';
import { useSelector } from 'react-redux';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';

const ActiveGames = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);

  // event handlers
  const onRefreshActiveGames = () => {
    dispatch(actions.refreshActiveGames(currentUser.id));
  };

  return <div>
       {currentUser.activeGameIDs ? currentUser.activeGameIDs.forEach((gID, index) => (
        <button
            key={gID}
            onClick={() => goToGame({gID})}
        >
        {!currentUser.activeGames || !currentUser.activeGames[index] ? 'no game or index?' : currentUser.activeGames[index]}
        </button>
      )) : 'NO GAMES.'}
      Your ({currentUser.name}) Active Games are: {currentUser.activeGameIDs}, {currentUser.activeGames}.
      <button
          onClick={onRefreshActiveGames}
          className='hover:text-white'
        >
          Refresh
        </button>
    </div>;
};

const goToGame = (gID) => {
    console.log(`Requesting game page for ID: ${gID}`)
}

export default ActiveGames;
