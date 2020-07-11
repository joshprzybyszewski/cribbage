import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import {
  sliceKey as authSliceKey,
  reducer as authReducer,
} from '../../../auth/slice';
import { gameSaga } from '../Game/saga';
import {
  sliceKey as gameSliceKey,
  reducer as gameReducer,
  actions as gameActions,
} from '../Game/slice';
import { selectActiveGames } from './selectors';
import { homeSaga } from './saga';
import {
  sliceKey as homeSliceKey,
  reducer as homeReducer,
  actions as homeActions,
} from './slice';

const ActiveGamesTable = () => {
  useInjectReducer({ key: authSliceKey, reducer: authReducer });
  useInjectSaga({ key: authSliceKey, saga: authSaga });
  useInjectReducer({ key: homeSliceKey, reducer: homeReducer });
  useInjectSaga({ key: homeSliceKey, saga: homeSaga });
  useInjectReducer({ key: gameSliceKey, reducer: gameReducer });
  useInjectSaga({ key: gameSliceKey, saga: gameSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGames = useSelector(selectActiveGames(currentUser.id));

  // event handlers
  const onRefreshActiveGames = () => {
    dispatch(homeActions.refreshActiveGames({ id: currentUser.id }));
  };
  const onGoToGame = gID => {
    dispatch(gameActions.goToGame(gID, history));
  };

  let gameButtons = activeGames.map(activeGame => {
    if (!activeGame || !activeGame.gameID) {
      return;
    }
    const gID = activeGame.gameID;

    return (
      <tr key={`gameRow ${gID}`}>
        <td className='active-games-table-data'>
          {activeGame.players
            .filter(p => p.id !== currentUser.id)
            .map(p => p.name)
            .join(', ')}
        </td>
        <td className='active-games-table-data active-games-table-data-sm'>
          {activeGame.players
            .filter(p => p.id === currentUser.id)
            .map(p => p.color)
            .toString()}
        </td>
        <td className='active-games-table-data active-games-table-data-sm'>
          {activeGame.created}
        </td>
        <td className='active-games-table-data active-games-table-data-sm'>
          {activeGame.lastMove}
        </td>
        <td className='active-games-table-data active-games-table-data-sm'>
          <button key={gID} onClick={() => onGoToGame(gID)}>
            Play!
          </button>
        </td>
      </tr>
    );
  });

  return (
    <div>
      <div className='flex flex-col'>
        <div className='-my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
          <div className='align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-200'>
            <table className='min-w-full'>
              <thead>
                <tr>
                  <th className='active-games-table-head active-games-table-head-text'>
                    Opponent(s)
                  </th>
                  <th className='active-games-table-head active-games-table-head-text'>
                    Your Color
                  </th>
                  <th className='active-games-table-head active-games-table-head-text'>
                    Started
                  </th>
                  <th className='active-games-table-head active-games-table-head-text'>
                    Last Move
                  </th>
                  <th className='active-games-table-head'>
                    <div
                      className='flex-shrink-0 h-5 w-5'
                      onClick={onRefreshActiveGames}
                    >
                      <img
                        className='h-10 w-10 rounded-full'
                        src='./refresh.svg'
                        alt='Refresh'
                      />
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white'>{gameButtons}</tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ActiveGamesTable;
