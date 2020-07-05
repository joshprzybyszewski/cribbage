import React from 'react';
import { useSelector } from 'react-redux';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser, selectActiveGames } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { gameSaga } from '../../../game/saga';
import { sliceKey as gameSliceKey, reducer as gameReducer, actions as gameActions } from '../../../game/slice';

const ActiveGames = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  useInjectReducer({ key: gameSliceKey, reducer: gameReducer });
  useInjectSaga({ key: gameSliceKey, saga: gameSaga });
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGames = useSelector(selectActiveGames);
  
  // event handlers
  const onRefreshActiveGames = () => {
    dispatch(actions.refreshActiveGames(currentUser.id));
  };
  const onGoToGame = (gID) => {
    dispatch(gameActions.goToGame(gID));
  };

  let gameButtons = [];
  if ( activeGames ) {
    for (const [gID, desc] of Object.entries(activeGames)) {
        if ( !gID || !desc ) {
            continue;
        }
        gameButtons.push(<br key={`br ${gameButtons.length}`}></br>)
        gameButtons.push(<button
            key={gID}
            onClick={onGoToGame(gID)}
        >
        {desc}
        </button>);
    }
  }

  return <div>
      This is supposed to be {currentUser.name}'s Active Games page.
      {gameButtons}
      <br></br>
      <button
          onClick={onRefreshActiveGames}
          className='hover:text-white'
        >
          Refresh
        </button>
    </div>;
};

const goToGame = (gID) => {
    // We're gonna need to navigate to the games page
    console.log(`Requesting game page for ID: ${gID}`)
}

export default ActiveGames;
