import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import ChoosingCard from './ChoosingCard';
import { TossSuggestion } from './slice';
import { useTossSuggestion } from './useTossSuggestion';


const SuggestionsTable = () => {
  const { suggestedHands } = useTossSuggestion();

  return (
    <TableContainer component={Paper}>
      <Table stickyHeader size='small' aria-label='Hand Suggestions table'>
        <TableHead>
          <TableRow>
            <TableCell>Hand (min)</TableCell>
            <TableCell>Hand (avg)</TableCell>
            <TableCell>Hand (max)</TableCell>
            <TableCell>Crib (min)</TableCell>
            <TableCell>Crib (avg)</TableCell>
            <TableCell>Crib (max)</TableCell>
            <TableCell>Toss</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {suggestedHands
            .filter((sug: TossSuggestion) => sug && sug.handPts && sug.cribPts && sug.toss)
            .map((sug: TossSuggestion) => {
              return (
                <TableRow hover>
                  <TableCell>
                    {sug.handPts.min}
                  </TableCell>
                  <TableCell>
                    {sug.handPts.avg}
                  </TableCell>
                  <TableCell>
                    {sug.handPts.max}
                  </TableCell>
                  <TableCell>
                    {sug.cribPts.min}
                  </TableCell>
                  <TableCell>
                    {sug.cribPts.avg}
                  </TableCell>
                  <TableCell>
                    {sug.cribPts.max}
                  </TableCell>
                  <TableCell>
                    <Grid
                      item
                      container
                      spacing={1}
                    >
                      <GridList>
                      {sug.toss.map((card) => (
                        <ChoosingCard
                          key={`tossCard${card.name}`}
                          card={card}
                          notEditable
                        />
                      ))}
                      </GridList>
                    </Grid>
                  </TableCell>
                </TableRow>
              );
            })}
        </TableBody>
      </Table>
    </TableContainer>
  )
};

export default SuggestionsTable;
