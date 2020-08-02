import React from 'react';

import Grid from '@material-ui/core/Grid';
import CountHandAction from 'app/containers/Game/CountHandAction';
import CribAction from 'app/containers/Game/CribAction';
import CutAction from 'app/containers/Game/CutAction';
import DealAction from 'app/containers/Game/DealAction';
import PegAction from 'app/containers/Game/PegAction';
import { gameSaga } from 'app/containers/Game/saga';
import { sliceKey, reducer } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const Action = ({ phase, isBlocking }) => {
  switch (phase) {
    case 'Deal':
      return <DealAction isBlocking={isBlocking} />;
    case 'BuildCrib':
      return <CribAction isBlocking={isBlocking} />;
    case 'Cut':
      return <CutAction isBlocking={isBlocking} />;
    case 'Pegging':
      return <PegAction isBlocking={isBlocking} />;
    case 'Counting':
      return <CountHandAction isBlocking={isBlocking} isCrib={false} />;
    case 'CribCounting':
      return <CountHandAction isBlocking={isBlocking} isCrib={true} />;
    default:
      return 'dev error!';
  }
};

const ActionBox = ({ phase, isBlocking }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  return (
    <Grid item container justify='center' spacing={1}>
      <Action phase={phase} isBlocking={isBlocking} />
    </Grid>
  );
};

ActionBox.propTypes = {
  phase: PropTypes.string.isRequired,
  isBlocking: PropTypes.bool.isRequired,
};

export default ActionBox;
