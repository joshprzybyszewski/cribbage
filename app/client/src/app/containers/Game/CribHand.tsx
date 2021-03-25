import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import { nanoid } from '@reduxjs/toolkit';

import { Card } from './models';
import PlayingCard from './PlayingCard';

interface Props {
    cards: Card[];
}

const CribHand: React.FunctionComponent<Props> = ({ cards }) => {
    return (
        !cards || (
            <Grid item container direction='row' justify='center' spacing={1}>
                <GridList>
                    {cards.map(card => (
                        <PlayingCard
                            disabled={false}
                            key={
                                card.name === 'unknown'
                                    ? `cribcard-unknown-${nanoid()}`
                                    : `cribcard-${card.name}`
                            }
                            card={card}
                        />
                    ))}
                </GridList>
            </Grid>
        )
    );
};

export default CribHand;
