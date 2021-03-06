import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import PlayingCard from 'app/containers/Game/PlayingCard';
import PropTypes from 'prop-types';

const showOpponentsHand = phase => {
    return phase !== 'Deal';
};

const PlayerHand = ({ hand, phase, side, pegged, mine }) => {
    if (!hand || !showOpponentsHand(phase)) {
        return null;
    }

    return (
        <Grid
            item
            container
            direction={side ? 'column' : 'row'}
            justify='center'
            spacing={1}
        >
            <GridList>
                {hand.map((card, index) => (
                    <PlayingCard
                        key={`handcard${index}`}
                        card={card}
                        mine={mine}
                        disabled={
                            phase === 'Pegging' &&
                            pegged &&
                            pegged.some(pc => pc.card.name === card.name)
                        }
                    />
                ))}
            </GridList>
        </Grid>
    );
};

PlayerHand.propTypes = {
    hand: PropTypes.array.isRequired,
    phase: PropTypes.string.isRequired,
    side: PropTypes.string.isRequired,
    pegged: PropTypes.array.isRequired,
    mine: PropTypes.bool.isRequired,
};

export default PlayerHand;
