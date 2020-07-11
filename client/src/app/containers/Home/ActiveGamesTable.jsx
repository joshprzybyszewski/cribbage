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

  let gameButtons = [];
  if (activeGames) {
    for (const [gID, activeGame] of Object.entries(activeGames)) {
      if (!gID || !activeGame) {
        continue;
      }

      let opponents = [];
      for (const [pID, pName] of Object.entries(activeGame.players)) {
        if (pID === currentUser.id) {
          continue;
        }
        opponents.push(pName);
      }

      gameButtons.push(
        <tr>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
            {opponents.join(', ')}
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            {activeGame.colors[currentUser.id]}
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            {activeGame.created}
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            {activeGame.lastMove}
          </td>
          <td class='px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 text-sm leading-5 font-medium'>
            <button key={gID} onClick={() => onGoToGame(gID)}>
              Play!
            </button>
          </td>
        </tr>,
      );
    }
  }

  return (
    <div>
      <div class='flex flex-col'>
        <div class='-my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
          <div class='align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-200'>
            <table class='min-w-full'>
              <thead>
                <tr>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Opponent(s)
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Your Color
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Started
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Last Move
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50'>
                    <div
                      class='flex-shrink-0 h-10 w-10'
                      onClick={onRefreshActiveGames}
                    >
                      <img
                        class='h-10 w-10 rounded-full'
                        src='./refresh.png'
                        alt='Refresh'
                      />
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody class='bg-white'>{gameButtons}</tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ActiveGamesTable;
