import React, { useState } from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Slider from '@material-ui/core/Slider';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import { sliceKey, reducer } from 'app/containers/Suggestions/slice';
import PropTypes from 'prop-types';
import { useInjectReducer } from 'redux-injectors';

const useStyles = makeStyles({
  root: {
    width: 120,
    height: 160,
    // display: 'flex',
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
  // fauxCardWrapper: {
  //   display: 'flex',
  //   flexDirection: 'column',
  // },
  // fauxCard: {
  //   flex: '1 0 auto',
  // },
  valueSlider: {
    // width: '5%',
    // height: '100%',
    width: '100%',
    height: '100%',
    position: 'relative',
    top: '0',
    left: '0',
  }
});

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

  const classes = useStyles();
  const value = getValue(card.toUpperCase());
  const [curVal, setCurVal] = useState(value);

  const suitVal = getSuitValue(card.toUpperCase());
  const [curSuit, setCurSuit] = useState(suitVal);


  const useRed = card.toUpperCase().includes('H') || card.toUpperCase().includes('D');

  const updateValue = () => {
    // TODO update the state so that this card increments suits
    console.log(`updateValue: ${notEditable} ${card}, ${curVal}`);
    setCurVal(prev => (prev % 13) + 1);
  }
  const updateSuit = (_) => {
    // TODO update the state so that this card increments suits
    console.log(`updateSuit: ${notEditable} ${card}, ${curSuit}`);
    setCurSuit(prev => (prev % 4) + 1);
  };

  return (
    <Card
      className={classes.root}>
      <div
        className={classes.fauxCardWrapper}
      >
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
      {
        !notEditable &&
        <div
          className={classes.valueSlider}
        >
          <CardContent>
            <Slider
              orientation='vertical'
              defaultValue={curVal}
              ValueLabelComponent={ValueLabelComponent}
              getAriaValueText={v => {
                console.log(`getAriaValueText: ${v}`);
                getValueString(v);
              }
              }
              aria-labelledby="discrete-slider"
              valueLabelDisplay="auto"
              step={1}
              marks
              min={1}
              max={13}
              onChange={event => {
                console.log(`onChange: ${event.target.value}`);
                // dispatch(actions.claimPoints(Number(event.target.value) / 100));
              }}
            />
          </CardContent>
        </div>
      }
      </div>
    </Card>
  );
};

ChoosingCard.propTypes = {
  card: PropTypes.string.isRequired,
  notEditable: PropTypes.bool,
};

export default ChoosingCard;
