import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import { nanoid } from '@reduxjs/toolkit';

import { Card, Phase } from './models';
import PlayingCard from './PlayingCard';
import { useGame } from './useGame';

const showOpponentsHand = (phase: Phase) => phase !== 'Deal';

interface Props {
    hand: Card[];
    side?: boolean;
    mine?: boolean;
}

const PlayerHand: React.FunctionComponent<Props> = ({ hand, side, mine }) => {
    const {
        game: { pegged_cards: peggedCards, phase },
    } = useGame();

    if (!hand || (!mine && !showOpponentsHand(phase))) {
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
                {hand.map(card => (
                    <PlayingCard
                        key={
                            card.name === 'unknown'
                                ? `handcard-unknown-${nanoid()}`
                                : `handcard-${card.name}`
                        }
                        card={card}
                        mine={mine}
                        disabled={
                            phase === 'Pegging' &&
                            peggedCards &&
                            peggedCards.some(pc => pc.card.name === card.name)
                        }
                    />
                ))}
            </GridList>
        </Grid>
    );
};

PlayerHand.defaultProps = {
    mine: false,
    side: false,
};

export default PlayerHand;
