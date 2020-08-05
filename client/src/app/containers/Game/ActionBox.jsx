import React from 'react';

import Grid from '@material-ui/core/Grid';
import makeStyles from '@material-ui/core/styles/makeStyles';
import CountHandAction from 'app/containers/Game/CountHandAction';
import CribAction from 'app/containers/Game/CribAction';
import CutAction from 'app/containers/Game/CutAction';
import DealAction from 'app/containers/Game/DealAction';
import PegAction from 'app/containers/Game/PegAction';
import PropTypes from 'prop-types';

const useStyles = makeStyles(theme => ({
  leftMargin: {
    marginLeft: theme.spacing(1),
  },
}));

const Action = ({ phase, isBlocking }) => {
  const classes = useStyles();

  switch (phase) {
    case 'Deal':
      return <DealAction isBlocking={isBlocking} styles={classes} />;
    case 'BuildCrib':
      return <CribAction isBlocking={isBlocking} />;
    case 'Cut':
      return <CutAction isBlocking={isBlocking} />;
    case 'Pegging':
      return <PegAction isBlocking={isBlocking} />;
    case 'Counting':
      return (
        <CountHandAction
          isBlocking={isBlocking}
          isCrib={false}
          styles={classes}
        />
      );
    case 'CribCounting':
      return (
        <CountHandAction
          isBlocking={isBlocking}
          isCrib={true}
          styles={classes}
        />
      );
    default:
      return 'dev error!';
  }
};

const ActionBox = ({ phase, isBlocking }) => {
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
