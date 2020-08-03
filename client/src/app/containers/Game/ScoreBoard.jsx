import React from 'react';

import Container from '@material-ui/core/Container';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import PlayerIcon from 'app/components/PlayerIcon';
import PropTypes from 'prop-types';

const ScoreBoard = ({ current_dealer, teams }) => {
  return (
    <Container fixed width='35px' size='small'>
      <TableContainer component={Paper} size='small'>
        <Table stickyHeader size='small' padding='none' aria-label='scoreboard'>
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
                    <PlayerIcon color={t.color} />
                  </TableCell>
                  <TableCell>{t.current_score}</TableCell>
                  <TableCell>{t.lag_score}</TableCell>
                  <TableCell>
                    {t.players.map(p => {
                      const s = p.name;
                      if (p.id === current_dealer) {
                        return (
                          <span key='dealer span'>
                            <strong>{p.name}</strong>
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

ScoreBoard.propTypes = {
  current_dealer: PropTypes.string.isRequired,
  teams: PropTypes.array.isRequired,
};

export default ScoreBoard;
