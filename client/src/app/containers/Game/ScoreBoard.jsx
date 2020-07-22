import React from 'react';

import blue from '@material-ui/core/colors/blue';
import green from '@material-ui/core/colors/green';
import grey from '@material-ui/core/colors/grey';
import red from '@material-ui/core/colors/red';
import Container from '@material-ui/core/Container';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';

// TODO import from a utils instead of redeclaring
// or even make an "atom" that is the icon
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
                    <PersonPinCircleIcon
                      style={{
                        color: myColorToHue(t.color),
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

export default ScoreBoard;
