import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
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
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGames = useSelector(selectActiveGames(currentUser.id));

  // event handlers
  const onRefreshActiveGames = () => {
    dispatch(homeActions.refreshActiveGames({ id: currentUser.id }));
  };

  let createData = (gameID, startedTime, lastMoveTime, opponents, color) => {
    let myColor = color
      ? color.includes('red')
        ? red[800]
        : color.includes('blue')
        ? blue[800]
        : color.includes('green')
        ? green[800]
        : grey[400]
      : grey[400];
    return { gameID, startedTime, lastMoveTime, opponents, myColor };
  };

  let rows = activeGames.map(activeGame => {
    if (!activeGame || !activeGame.gameID) {
      return;
    }

    return createData(
      activeGame.gameID,
      activeGame.created,
      activeGame.lastMove,
      activeGame.players
        .filter(p => p.id !== currentUser.id)
        .map(p => p.name)
        .join(', '),
      activeGame.players.filter(p => p.id === currentUser.id).map(p => p.color),
    );
  });

  return (
    <TableContainer component={Paper}>
      <Table size='small' aria-label='active games table'>
        <TableHead>
          <TableRow>
            <TableCell>Opponent(s)</TableCell>
            <TableCell>Your Color</TableCell>
            <TableCell>Started</TableCell>
            <TableCell>Last Move</TableCell>
            <TableCell>
              <IconButton aria-label='refresh' onClick={onRefreshActiveGames}>
                <RefreshIcon />
              </IconButton>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map(row => (
            <TableRow key={row.name}>
              <TableCell component='th' scope='row'>
                {row.opponents}
              </TableCell>
              <TableCell>
                <PersonPinCircleIcon style={{ color: row.myColor }} />
              </TableCell>
              <TableCell>{row.startedTime}</TableCell>
              <TableCell>{row.lastMoveTime}</TableCell>
              <TableCell>
                <IconButton
                  aria-label='play'
                  onClick={() => goToGame(row.gameID)}
                >
                  <SportsEsportsIcon />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

const goToGame = gID => {
  // We're gonna need to navigate to the games page
  console.log(`Will request game page for ID: ${gID}`);
};

export default ActiveGamesTable;
