import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import Slider from '@material-ui/core/Slider';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import { gameSaga } from 'app/containers/Game/saga';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const CutAction = ({ isBlocking }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const [sliderVal, setSliderVal] = useState(0.5);

  const dispatch = useDispatch();

  return (
    <div>
      <Slider
        value={sliderVal}
        disabled={!isBlocking}
        orientation='vertical'
        getAriaValueText={value => {
          return `${value}%`;
        }}
        defaultValue={50}
        aria-labelledby='vertical-slider'
        onChange={(_, newValue) => {
          setSliderVal(newValue);
        }}
      />
      <Button
        disabled={!isBlocking}
        variant='contained'
        color='primary'
        endIcon={<CallSplitIcon />}
        onClick={() => {
          dispatch(actions.cutDeck(sliderVal / 100));
        }}
      >
        Cut
      </Button>
    </div>
  );
};

CutAction.propTypes = {
  isBlocking: PropTypes.bool.isRequired,
};

export default CutAction;
