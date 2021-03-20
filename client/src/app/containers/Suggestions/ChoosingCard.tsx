import React from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Slider from '@material-ui/core/Slider';
import { makeStyles } from '@material-ui/core/styles';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import clsx from 'clsx';

import { 
  Card as ModelCard,
  Value as ModelValue,
  Suit as ModelSuit,
} from '../Game/models';
import { useTossSuggestion } from './useTossSuggestion';

function getValueString(val: number): string {
  if (val === 13) {
    return 'K';
  }
  if (val === 12) {
    return 'Q';
  }
  if (val === 11) {
    return 'J';
  }
  if (val === 1) {
    return 'A';
  }

  return `${val}`;
}

function getCard(val: ModelValue, s: ModelSuit): ModelCard {
  return {
    name: `${getValueString(val)}${s[0]}`,
    suit: s,
    value: val,
  }
}

function getUpdatedSuit(card: ModelCard): ModelCard {
  switch (card.suit) {
    case 'Spades':
       return getCard(card.value, 'Clubs');
      case 'Clubs':
       return getCard(card.value, 'Diamonds');
      case 'Diamonds':
       return getCard(card.value, 'Hearts');
      case 'Hearts':
       return getCard(card.value, 'Spades');
       default:
  return getCard(1, 'Spades');
  }
}

function getSuitEmoji(card: ModelCard): string {
  switch (card.suit) {
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
}

function getUpdatedValue(card: ModelCard, val: ModelValue): ModelCard {
  return getCard(val, card.suit);
}


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
  fauxCard: {},
  fauxCardWrapper: {
    // flexGrow: '1',
    height: '100%',
  },
  valueSlider: {
    flexBasis: '10%',
    height: '50%',
    marginTop: '10%',
  },
  redCard: {
    color: red[700],
},
blackCard: {
    color: 'black',
}
});


const ValueLabelComponent: React.FunctionComponent = (props) => {
  // function ValueLabelComponent(props) {
  const { children, open, value } = props;

  return (
    <Tooltip open={open} enterTouchDelay={0} placement="top" title={getValueString(value)}>
      {children}
    </Tooltip>
  );
}

interface Props {
  card: ModelCard;
  notEditable?: boolean;
}

const ChoosingCard: React.FunctionComponent<Props> = ({ card, notEditable }) => {
  const { updateCard } = useTossSuggestion();
  
  const classes = useStyles();

  const useRed = !['Spades', 'Clubs'].includes(card.suit);

  const updateValue = notEditable ?
  () => {} :
  (v: number) => updateCard(
    card,
    getUpdatedValue(card, v as ModelValue),
  );

  const updateSuit = notEditable ?
  () => {} :
  () => updateCard(
        card,
       getUpdatedSuit(card),
      );

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
              className={clsx(classes.value, {
                [classes.redCard]: useRed,
                [classes.blackCard]: !useRed,
              })}
              gutterBottom
            >
              {card.value}
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
            defaultValue={card.value}
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

export default ChoosingCard;
