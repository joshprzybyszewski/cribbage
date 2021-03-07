/* eslint-disable jsx-a11y/no-static-element-interactions */
/* eslint-disable jsx-a11y/click-events-have-key-events */
// TODO don't disable eslint. maybe use a button instead
import React from 'react';

import {
    Card as CardComponent,
    CardContent,
    Typography,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';

import { Card } from './models';
import { useGame } from './useGame';

const useStyles = makeStyles({
    root: {
        width: 120,
        height: 160,
    },
    value: {
        fontSize: 14,
    },
    red: {
        color: 'red',
    },
    black: {
        color: 'black',
    },
    suit: {
        justifyContent: 'center',
        alignItems: 'center',
        verticalAlign: 'center',
        textAlign: 'center',
    },
});

const mapSuitToSymbol = (
    suit: 'Spades' | 'Clubs' | 'Diamonds' | 'Hearts',
): string => {
    switch (suit) {
        case 'Spades':
            return '♠️';
        case 'Clubs':
            return '♣️';
        case 'Diamonds':
            return '♦️';
        case 'Hearts':
            return '♥️';
        default:
            return '?';
    }
};

interface Props {
    card: Card;
    disabled: boolean;
    experimental?: boolean;
    mine?: boolean;
}

const PlayingCard: React.FunctionComponent<Props> = ({
    card,
    disabled,
    experimental,
    mine,
}) => {
    const { selectedCards, toggleSelectedCard } = useGame();
    const classes = useStyles();
    const useRed = !['Spades', 'Clubs'].includes(card.suit);

    if (experimental) {
        return (
            <CardComponent className={classes.root}>
                <CardContent>
                    <Typography
                        className={clsx(classes.value, {
                            [classes.black]: !useRed,
                            [classes.red]: useRed,
                        })}
                        gutterBottom
                    >
                        {card.value}
                    </Typography>
                    <Typography className={classes.suit}>
                        {mapSuitToSymbol(card.suit)}
                    </Typography>
                </CardContent>
            </CardComponent>
        );
    }

    if (!card) {
        return null;
    }
    if (card.name === 'unknown') {
        // Currently, this returns a grayed out box, but it should show
        // a back of a card
        return (
            <div className='w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800' />
        );
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
            className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
                disabled ? 'bg-gray-500' : 'bg-white'
            } ${useRed ? 'text-red-700' : 'text-black'}`}
            style={{
                position: 'relative',
                top: chosen ? '-10px' : '',
            }}
        >
            {card.name}
        </div>
    );
};

PlayingCard.defaultProps = {
    experimental: false,
    mine: false,
};

export default PlayingCard;
