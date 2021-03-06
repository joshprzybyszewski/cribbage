import React from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const useStyles = makeStyles({
    root: {
        width: 120,
        height: 160,
    },
    value: {
        fontSize: 14,
    },
    suit: {
        justifyContent: 'center',
        alignItems: 'center',
        verticalAlign: 'center',
        textAlign: 'center',
    },
});

const PlayingCard = ({ card, disabled, experimental, mine }) => {
    useInjectReducer({ key: sliceKey, reducer: reducer });
    useInjectSaga({ key: sliceKey, saga: gameSaga });
    const classes = useStyles();
    const dispatch = useDispatch();
    const currentAction = useSelector(selectCurrentAction);

    const useRed = !['Spades', 'Clubs'].includes(card.suit);

    if (experimental) {
        return (
            <Card className={classes.root}>
                <CardContent>
                    <Typography
                        className={classes.value}
                        color={useRed ? 'red' : 'black'}
                        gutterBottom
                    >
                        {card.value}
                    </Typography>
                    <Typography className={classes.suit}>
                        {card.suit === 'Spades'
                            ? '♠️'
                            : card.suit === 'Clubs'
                            ? '♣️'
                            : card.suit === 'Diamonds'
                            ? '♦️'
                            : card.suit === 'Hearts'
                            ? '♥️'
                            : '?'}
                    </Typography>
                </CardContent>
            </Card>
        );
    }

    if (!card) {
        return null;
    } else if (card.name === 'unknown') {
        // Currently, this returns a grayed out box, but it should show
        // a back of a card
        return (
            <div className='w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800' />
        );
    }

    const chosen = currentAction.selectedCards.indexOf(card) !== -1;
    const toggleChosen = () => {
        if (!disabled) {
            dispatch(actions.selectCard(card));
        }
    };

    return (
        <div
            onClick={mine ? toggleChosen : () => {}}
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

PlayingCard.propTypes = {
    card: PropTypes.object.isRequired,
    disabled: PropTypes.bool.isRequired,
    experimental: PropTypes.bool.isRequired,
    mine: PropTypes.bool.isRequired,
};

export default PlayingCard;
