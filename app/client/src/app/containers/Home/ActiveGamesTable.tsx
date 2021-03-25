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
import { useHistory } from 'react-router-dom';

import { useAuth } from '../../../auth/useAuth';
import { useGame } from '../Game/useGame';
import ActiveGames from './ActiveGames';
import { useActiveGames } from './useActiveGames';

const ActiveGamesTable = () => {
    const { games, refreshGames } = useActiveGames();
    const { loadActiveGame } = useGame();
    const { currentUser } = useAuth();
    const history = useHistory();

    const onGoToGame = async (id: number) => {
        await loadActiveGame(id);
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
