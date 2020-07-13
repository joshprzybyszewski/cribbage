import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import blue from '@material-ui/core/colors/blue';
import grey from '@material-ui/core/colors/grey';
import green from '@material-ui/core/colors/green';
import red from '@material-ui/core/colors/red';
import IconButton from '@material-ui/core/IconButton';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';
import RefreshIcon from '@material-ui/icons/Refresh';
import SportsEsportsIcon from '@material-ui/icons/SportsEsports';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

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

const myColorToHue = color => {
  return color
    ? color.includes('red')
      ? red[800]
      : color.includes('blue')
      ? blue[800]
      : color.includes('green')
      ? green[800]
      : grey[400]
    : grey[400];
};

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

  return (
    <TableContainer
      component={Paper}
      style={{
        maxHeight: 500,
      }}
    >
      <Table stickyHeader size='small' aria-label='active games table'>
        <TableHead>
          <TableRow>
            <TableCell>Other Player(s)</TableCell>
            <TableCell>Your Color</TableCell>
            <TableCell>Started</TableCell>
            <TableCell>Last Activity</TableCell>
            <TableCell>
              <IconButton aria-label='refresh' onClick={onRefreshActiveGames}>
                <RefreshIcon />
              </IconButton>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {activeGames
            .filter(ag => ag && ag.gameID)
            .map(ag => {
              return (
                <TableRow hover key={ag.gameID}>
                  <TableCell component='th' scope='row'>
                    {ag.players
                      .filter(p => p.id !== currentUser.id)
                      .map(p => p.name)
                      .join(', ')}
                  </TableCell>
                  <TableCell>
                    <PersonPinCircleIcon
                      style={{
                        color: myColorToHue(
                          ag.players
                            .filter(p => p.id === currentUser.id)
                            .map(p => p.color),
                        ),
                      }}
                    />
                  </TableCell>
                  <TableCell>{ag.created}</TableCell>
                  <TableCell>{ag.lastMove}</TableCell>
                  <TableCell>
                    <IconButton
                      aria-label='play'
                      onClick={() => onGoToGame(ag.gameID)}
                    >
                      <SportsEsportsIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default ActiveGamesTable;
