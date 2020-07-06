import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { selectActiveGames } from '../../../home/selectors';
import { homeSaga } from '../../../home/saga';
import { sliceKey as homeSliceKey, reducer as homeReducer, actions as homeActions } from '../../../home/slice';

const ActiveGamesTable = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  useInjectReducer({ key: homeSliceKey, reducer: homeReducer });
  useInjectSaga({ key: homeSliceKey, saga: homeSaga });
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGames = useSelector(selectActiveGames(currentUser.id));
  
  // event handlers
  const onRefreshActiveGames = () => {
    dispatch(homeActions.refreshActiveGames({ id: currentUser.id }));
  };

  let gameButtons = [];
  if ( activeGames ) {
    for (const [gID, desc] of Object.entries(activeGames)) {
        if ( !gID || !desc ) {
            continue;
        }
        gameButtons.push(<br key={`br ${gID}`}></br>)
        gameButtons.push(<button
            key={gID}
            onClick={() => goToGame({gID})}
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

export default ActiveGamesTable;
