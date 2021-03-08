import React, { useState } from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Slider from '@material-ui/core/Slider';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import {
  sliceKey,
  reducer,
  actions as sugActions,
} from 'app/containers/Suggestions/slice';
import PropTypes from 'prop-types';
import { useInjectReducer } from 'redux-injectors';
import { useDispatch } from 'react-redux';

const useStyles = makeStyles({
  root: {
    width: 120,
    height: 160,
    display: 'flex',
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
  fauxCardWrapper: {
    flexGrow: '1',
    height: '100%',
  },
  valueSlider: {
    flexBasis: '10%',
    height: '50%',
    marginTop: '10%',
  }
});


function ValueLabelComponent(props) {
  const { children, open, value } = props;

  return (
    <Tooltip open={open} enterTouchDelay={0} placement="top" title={getValueString(value)}>
      {children}
    </Tooltip>
  );
}

const ChoosingCard = ({ card, notEditable }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  const dispatch = useDispatch();

  const classes = useStyles();
  const value = getValue(card.toUpperCase());
  const suitVal = getSuitValue(card.toUpperCase());

  const useRed = suitVal > 2;

  const updateValue = (v) => {
    !notEditable &&
      dispatch(sugActions.updateCard({
        card: card,
        newCard: getUpdatedValue(card, v),
      }));
  }
  const updateSuit = (_) => {
    !notEditable &&
      dispatch(sugActions.updateCard({
        card: card,
        newCard: getUpdatedSuit(card),
      }));
  };

  return (
    <div
      className={classes.root}
    >
      <div
        className={classes.fauxCardWrapper}
      >
        <Card>
          <CardContent
            className={classes.fauxCard}
            onClick={updateSuit}
          >
            <Typography
              className={classes.value}
              color={useRed ? 'red' : 'black'}
              gutterBottom
            >
              {getValueString(value)}
            </Typography>
            <Typography className={classes.suit}>
              {getSuitEmoji(card)}
            </Typography>
          </CardContent>
        </Card>
      </div>
      {
        !notEditable &&
        <div
          className={classes.valueSlider}
        >
          <Slider
            orientation='vertical'
            defaultValue={value}
            ValueLabelComponent={ValueLabelComponent}
            getAriaValueText={v => getValueString(v)}
            aria-labelledby="discrete-value-slider"
            valueLabelDisplay="auto"
            step={1}
            marks
            min={1}
            max={13}
            onChangeCommitted={(_, v) => updateValue(v)}
          />
        </div>
      }
    </div>
  );
};

ChoosingCard.propTypes = {
  card: PropTypes.string.isRequired,
  notEditable: PropTypes.bool,
};

export default ChoosingCard;

function getSuitValue(card) {
  if (card.includes('C')) {
    return 2;
  } else if (card.includes('S')) {
    return 1;
  } else if (card.includes('H')) {
    return 4;
  } else if (card.includes('D')) {
    return 3;
  }

  return 0;
}

function getUpdatedSuit(card) {
  if (card.includes('S')) {
    return card.replace('S', 'C');
  } else if (card.includes('C')) {
    return card.replace('C', 'D');
  } else if (card.includes('D')) {
    return card.replace('D', 'H');
  }
  return card.replace('H', 'S');
}

function getSuitEmoji(card) {
  if (card.includes('C')) {
    return '♣️';
  } else if (card.includes('S')) {
    return '♠️';
  } else if (card.includes('H')) {
    return '♥️';
  } else if (card.includes('D')) {
    return '♦️';
  }

  return '?';
}

function getValue(card) {
  if (card.startsWith('K') || card.startsWith('13')) {
    return 13; // 'K';
  } else if (card.startsWith('Q') || card.startsWith('12')) {
    return 12; // 'Q';
  } else if (card.startsWith('J') || card.startsWith('11')) {
    return 11; // 'J';
  } else if (card.startsWith('10')) {
    return 10; // '10';
  } else if (card.startsWith('A') || card.startsWith('1')) {
    return 1; // 'A';
  }

  return parseInt(card.substr(0, card.length - 1));
}

function getValueString(val) {
  if (val === 13) {
    return 'K';
  } else if (val === 12) {
    return 'Q';
  } else if (val === 11) {
    return 'J';
  } else if (val === 1) {
    return 'A';
  }

  return `${val}`;
}

function getUpdatedValue(card, val) {
  if (card.includes('S')) {
    return getValueString(val) + 'S';
  } else if (card.includes('C')) {
    return getValueString(val) + 'C';
  } else if (card.includes('D')) {
    return getValueString(val) + 'D';
  }

  return getValueString(val) + 'H';
}