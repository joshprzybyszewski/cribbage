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
import ChoosingCard from 'app/containers/Suggestions/ChoosingCard';

import { selectSuggestions } from 'app/containers/Suggestions/selectors';
import { sliceKey, reducer } from 'app/containers/Suggestions/slice';
import { useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const SuggestionsTable = () => {
  const sugs = useSelector(selectSuggestions);

  if (!Array.isArray(sugs)) {
    return (`suggestions are: ${sugs}`);
  }

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
          {sugs
            .filter(sug => sug && sug.handPts && sug.cribPts && sug.toss)
            .map(sug => {
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
                      {sug.toss.map((card, index) => (
                        <ChoosingCard
                          key={`tossCard${index}`}
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
