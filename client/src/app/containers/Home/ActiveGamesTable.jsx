import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { selectActiveGames } from '../../../home/selectors';
import { homeSaga } from '../../../home/saga';
import {
  sliceKey as homeSliceKey,
  reducer as homeReducer,
  actions as homeActions,
} from '../../../home/slice';

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
  if (activeGames) {
    for (const [gID, desc] of Object.entries(activeGames)) {
      if (!gID || !desc) {
        continue;
      }

      gameButtons.push(
        <tr>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
            {desc}
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            (not implemented)
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            (not implemented)
          </td>
          <td class='px-6 py-4 whitespace-no-wrap border-b border-gray-200 text-sm leading-5 text-gray-500'>
            (not implemented)
          </td>
          <td class='px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 text-sm leading-5 font-medium'>
            <button key={gID} onClick={() => goToGame({ gID })}>
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
                    Opponents
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Your Color
                  </th>
                  <th class='px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider'>
                    Created
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

const goToGame = gID => {
  // We're gonna need to navigate to the games page
  console.log(`Requesting game page for ID: ${gID}`);
};

export default ActiveGamesTable;
