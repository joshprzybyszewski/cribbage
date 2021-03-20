/* eslint-disable jsx-a11y/no-static-element-interactions */
/* eslint-disable jsx-a11y/click-events-have-key-events */
// TODO don't disable eslint. maybe use a button instead
import React from 'react';

import { makeStyles, Typography } from '@material-ui/core';
import { grey, red } from '@material-ui/core/colors';
import clsx from 'clsx';

import { ReactComponent as CardBack } from './card-back.svg';
import { Card } from './models';
import { useGame } from './useGame';

const useStyles = makeStyles(theme => ({
    cardSize: {
        width: theme.spacing(6),
        height: theme.spacing(8),
    },
    cardBase: {
        textAlign: 'center',
        display: 'inline-block',
        borderStyle: 'solid',
        borderWidth: 2,
        borderColor: 'black',
        position: 'relative',
    },
    disabledCard: {
        backgroundColor: grey[500],
    },
    redCard: {
        color: red[700],
    },
    blackCard: {
        color: 'black',
    },
    selected: {
        top: -theme.spacing(1),
    },
}));

interface Props {
    card: Card;
    disabled: boolean;
    mine?: boolean;
}

const PlayingCard: React.FunctionComponent<Props> = ({
    card,
    disabled,
    mine,
}) => {
    const classes = useStyles();
    const { selectedCards, toggleSelectedCard } = useGame();
    const useRed = !['Spades', 'Clubs'].includes(card.suit);

    if (!card) {
        return null;
    }
    if (card.name === 'unknown') {
        return <CardBack className={classes.cardSize} />;
    }

    const chosen = selectedCards.indexOf(card) !== -1;
    const handleClick = () => {
        if (!disabled && mine) {
            toggleSelectedCard(card);
        }
    };

    return (
        <div
            onClick={handleClick}
            className={clsx(classes.cardBase, classes.cardSize, {
                [classes.disabledCard]: disabled,
                [classes.redCard]: useRed,
                [classes.blackCard]: !useRed,
                [classes.selected]: chosen,
            })}
        >
            <Typography variant='body1'>{card.name}</Typography>
        </div>
    );
};

PlayingCard.defaultProps = {
    mine: false,
};

export default PlayingCard;
