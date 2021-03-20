import React from 'react';

import { IconButton, TableCell, TableRow } from '@material-ui/core';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';
import SportsEsportsIcon from '@material-ui/icons/SportsEsports';

import { User } from '../../../auth/slice';
import { colorToHue } from '../../../utils/colorToHue';
import { ActiveGame } from './slice';

const getHueForPlayerInGame = (player: User, game: ActiveGame) => {
    const color = game.players.find(p => p.id === player.id)?.color;
    return colorToHue(color ?? '');
};

interface Props {
    player: User;
    games: ActiveGame[];
    onClickPlay: (game: ActiveGame) => void;
}

const ActiveGames: React.FunctionComponent<Props> = ({
    player,
    games,
    onClickPlay,
}) => {
    return (
        <>
            {games
                .filter(ag => ag && ag.gameID)
                .map(ag => {
                    return (
                        <TableRow hover key={ag.gameID}>
                            <TableCell component='th' scope='row'>
                                {ag.players
                                    .filter(p => p.id !== player.id)
                                    .map(p => p.name)
                                    .join(', ')}
                            </TableCell>
                            <TableCell>
                                <PersonPinCircleIcon
                                    // TODO we shouldn't use inline styles, but for now it's aight
                                    style={{
                                        color: getHueForPlayerInGame(
                                            player,
                                            ag,
                                        ),
                                    }}
                                />
                            </TableCell>
                            <TableCell>{ag.created}</TableCell>
                            <TableCell>{ag.lastMove}</TableCell>
                            <TableCell>
                                <IconButton
                                    aria-label='play'
                                    onClick={() => onClickPlay(ag)}
                                >
                                    <SportsEsportsIcon />
                                </IconButton>
                            </TableCell>
                        </TableRow>
                    );
                })}
        </>
    );
};

export default ActiveGames;
