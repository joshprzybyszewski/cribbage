import React from 'react';

import {
    Container,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
} from '@material-ui/core';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';

import { colorToHue } from '../../../utils/colorToHue';
import { Team } from './models';

// TODO abstract network models from the models we use so we can fix naming
interface Props {
    current_dealer: string;
    teams: Team[];
}

const ScoreBoard: React.FunctionComponent<Props> = ({
    current_dealer,
    teams,
}) => {
    return (
        <Container fixed maxWidth='xs'>
            <TableContainer component={Paper}>
                <Table
                    stickyHeader
                    size='small'
                    padding='none'
                    aria-label='scoreboard'
                >
                    <TableHead>
                        <TableRow>
                            <TableCell>Color</TableCell>
                            <TableCell>Score</TableCell>
                            <TableCell>Lag</TableCell>
                            <TableCell>Names</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {teams.map(t => {
                            return (
                                <TableRow key={t.color}>
                                    <TableCell>
                                        <PersonPinCircleIcon
                                            style={{
                                                color: colorToHue(t.color),
                                            }}
                                        />
                                    </TableCell>
                                    <TableCell>{t.current_score}</TableCell>
                                    <TableCell>{t.lag_score}</TableCell>
                                    <TableCell>
                                        {t.players.map(p => {
                                            const s = p.name;
                                            if (p.id === current_dealer) {
                                                return (
                                                    <span key='dealer span'>
                                                        <strong>
                                                            {p.name}
                                                        </strong>
                                                        {' (dealer)'}
                                                    </span>
                                                );
                                            }
                                            return s;
                                        })}
                                    </TableCell>
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>
            </TableContainer>
        </Container>
    );
};

export default ScoreBoard;
