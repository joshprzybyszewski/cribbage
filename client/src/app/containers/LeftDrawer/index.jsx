import React from 'react';

import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import IconButton from '@material-ui/core/IconButton';
import makeStyles from '@material-ui/core/styles/makeStyles';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import { drawerWidth } from 'app/containers/LeftDrawer/constants';
import LoggedInDrawer from 'app/containers/LeftDrawer/LoggedInDrawer';
import LoggedOutDrawer from 'app/containers/LeftDrawer/LoggedOutDrawer';
import { selectLoggedIn } from 'auth/selectors';
import { sliceKey, reducer } from 'auth/slice';
import PropTypes from 'prop-types';
import { useSelector } from 'react-redux';
import { useInjectReducer } from 'redux-injectors';

const useStyles = makeStyles(theme => ({
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
  },
}));

const LeftDrawer = ({ isOpen, handleDrawerClose, mouseOverBinding }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  const loggedIn = useSelector(selectLoggedIn);

  const classes = useStyles();

  return (
    <Drawer
      className={classes.drawer}
      variant='persistent'
      anchor='left'
      open={isOpen}
      classes={{
        paper: classes.drawerPaper,
      }}
      {...mouseOverBinding}
    >
      <div className={classes.drawerHeader}>
        <IconButton onClick={handleDrawerClose}>
          <ChevronLeftIcon />
        </IconButton>
      </div>
      <Divider />
      {loggedIn ? <LoggedInDrawer /> : <LoggedOutDrawer />}
    </Drawer>
  );
};

LeftDrawer.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  handleDrawerClose: PropTypes.func.isRequired,
  mouseOverBinding: PropTypes.object.isRequired,
};

export default LeftDrawer;
