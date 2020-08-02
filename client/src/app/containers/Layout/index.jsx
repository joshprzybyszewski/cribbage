import React, { useState, useEffect } from 'react';

import makeStyles from '@material-ui/core/styles/makeStyles';
import Navbar from 'app/components/Navbar';
import Alert from 'app/containers/Alert';
import LeftDrawer from 'app/containers/LeftDrawer';
import {
  drawerWidth,
  drawerCloseDelay,
} from 'app/containers/LeftDrawer/constants';
import { selectLoggedIn } from 'auth/selectors';
import { sliceKey, reducer } from 'auth/slice';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { useSelector } from 'react-redux';
import { useInjectReducer } from 'redux-injectors';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(3),
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: -drawerWidth,
  },
  contentShift: {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: 0,
  },
}));

const Layout = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  const loggedIn = useSelector(selectLoggedIn);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const [mouseInDrawer, drawerMouseOverBinding] = useMouseIn();

  useEffect(() => {
    if (isDrawerOpen && !mouseInDrawer) {
      const timer = setTimeout(() => {
        setIsDrawerOpen(false);
      }, drawerCloseDelay);
      return () => clearTimeout(timer);
    }
  }, [isDrawerOpen, mouseInDrawer]);

  const handleDrawerOpen = () => {
    setIsDrawerOpen(true);
  };
  const handleDrawerClose = () => {
    setIsDrawerOpen(false);
  };

  const classes = useStyles();

  return (
    <React.Fragment>
      <Navbar loggedIn={loggedIn} handleDrawerOpen={handleDrawerOpen} />
      <div className={classes.root}>
        <LeftDrawer
          width={drawerWidth}
          mouseOverBinding={drawerMouseOverBinding}
          isOpen={isDrawerOpen}
          handleDrawerClose={handleDrawerClose}
        />
        <main
          className={clsx(classes.content, {
            [classes.contentShift]: isDrawerOpen,
          })}
        >
          <div className={classes.drawerHeader} />
          <Alert />
          {props.children}
        </main>
      </div>
    </React.Fragment>
  );
};

Layout.propTypes = {
  children: PropTypes.node.isRequired,
};

const useMouseIn = () => {
  const [mouseIn, setMouseIn] = useState(false);
  const onMouseEnter = () => {
    setMouseIn(true);
  };
  const onMouseLeave = () => {
    setMouseIn(false);
  };
  return [mouseIn, { onMouseEnter, onMouseLeave }];
};

export default Layout;
