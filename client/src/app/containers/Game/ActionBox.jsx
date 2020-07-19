import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Slider from '@material-ui/core/Slider';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentAction, selectCurrentGame } from './selectors';
import { DealAction } from './DealAction';
import { CribAction } from './CribAction';
import { CutAction } from './CutAction';
import { PegAction } from './PegAction';
import { CountHandAction } from './CountHandAction';
import { CountCribAction } from './CountCribAction';

const ActionBox = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  return (
    <Grid item container justify='center' spacing={1}>
      {props.phase === 'Deal' ? (
        <DealAction />
      ) : props.phase === 'BuildCrib' ? (
        <CribAction />
      ) : props.phase === 'Cut' ? (
        <CutAction />
      ) : props.phase === 'Pegging' ? (
        <PegAction />
      ) : props.phase === 'Counting' ? (
        <CountHandAction />
      ) : props.phase === 'CribCounting' ? (
        <CountCribAction />
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
