import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';

import PlayingCard from './PlayingCard';
import { Card } from './slice';

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
                            mine={false}
                            key={`cribcard-${card.name}`}
                            card={card}
                        />
                    ))}
                </GridList>
            </Grid>
        )
    );
};

export default CribHand;
