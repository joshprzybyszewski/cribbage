import React from 'react';

import {
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
} from '@material-ui/core';
import RefreshIcon from '@material-ui/icons/Refresh';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { useAuth } from '../../../auth/useAuth';
import { actions as gameActions } from '../Game/slice';
import ActiveGames from './ActiveGames';
import { useActiveGames } from './useActiveGames';

const ActiveGamesTable = () => {
    const { games, refreshGames } = useActiveGames();
    const { currentUser } = useAuth();
    const dispatch = useDispatch();
    const history = useHistory();

    const onGoToGame = (id: string) => {
        dispatch(gameActions.setGameID(id));
        history.push('/game');
    };

    return (
        <TableContainer component={Paper}>
            <Table stickyHeader size='small' aria-label='active games table'>
                <TableHead>
                    <TableRow>
                        <TableCell>Other Player(s)</TableCell>
                        <TableCell>Your Color</TableCell>
                        <TableCell>Started</TableCell>
                        <TableCell>Last Activity</TableCell>
                        <TableCell>
                            <IconButton
                                aria-label='refresh'
                                onClick={refreshGames}
                            >
                                <RefreshIcon />
                            </IconButton>
                        </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    <ActiveGames
                        player={currentUser}
                        games={games}
                        onClickPlay={g => onGoToGame(g.gameID)}
                    />
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default ActiveGamesTable;
