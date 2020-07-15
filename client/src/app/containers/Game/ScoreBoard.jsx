import React from 'react';

import blue from '@material-ui/core/colors/blue';
import grey from '@material-ui/core/colors/grey';
import green from '@material-ui/core/colors/green';
import red from '@material-ui/core/colors/red';

import Container from '@material-ui/core/Container';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

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
                      let s = p.name;
                      if (p.id === current_dealer) {
                        return (
                          <span>
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
  // return (
  //   <div className='py-2 px-4 bg-gray-400 border-gray-700 rounded-lg border-2'>
  //     Scores:{' '}
  //     {teams.map(t => (
  //       <div className='flex flex-row'>
  //         <div className='flex-1 capitalize'>{t.color}</div>
  //         <div className='flex-2'>{t.current_score}</div>
  //         <div className='flex-2'>
  //           ({t.players.map(p => p.name).join(', ')})
  //         </div>
  //       </div>
  //     ))}
  //   </div>
  // );
};

export default ScoreBoard;
