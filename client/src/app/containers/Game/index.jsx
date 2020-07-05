import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { gameSaga } from '../../../game/saga';
import { sliceKey, reducer, actions } from '../../../game/slice';
import { selectCurrentGame } from '../../../game/selectors';

const Game = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const activeGame = useSelector(selectCurrentGame);
  const activeGameID = activeGame.id;
  
  // event handlers
  const onRefreshCurrentGame = () => {
    dispatch(actions.refreshGame(activeGameID, history));
  };

  let gameResp = [];
  if ( activeGame ) {
    for (const [key, val] of Object.entries(activeGame)) {
        gameResp.push(`${key}: ${val} `);
        gameResp.push(<br key={`br ${key}`}></br>);
    }
  }


  return <div>
      This will be a page for the game of a user.
      <br></br>
      {gameResp}
      <br></br>
      <button
          onClick={onRefreshCurrentGame}
          className='hover:text-white'
        >
          Refresh
        </button>
    </div>;
};

export default Game;
