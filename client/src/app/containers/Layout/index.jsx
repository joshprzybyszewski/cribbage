import React, { useState, useRef } from 'react';

import makeStyles from '@material-ui/core/styles/makeStyles';
import Navbar from 'app/components/Navbar';
import Alert from 'app/containers/Alert';
import LeftDrawer from 'app/containers/LeftDrawer';
import { selectLoggedIn } from 'auth/selectors';
import { sliceKey, reducer } from 'auth/slice';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { useSelector } from 'react-redux';
import { useInjectReducer } from 'redux-injectors';

const drawerWidth = 240;

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
  const [drawerOpen, setDrawerOpen] = useState(false);
  // use a ref so we can reset the timer
  const timerID = useRef(-1);

  const handleDrawerOpen = () => {
    setDrawerOpen(true);
    timerID.current = setTimeout(() => {
      setDrawerOpen(false);
    }, 5000);
  };
  const handleDrawerClose = () => {
    setDrawerOpen(false);
    clearTimeout(timerID.current);
    timerID.current = -1;
  };

  const classes = useStyles();

  return (
    <React.Fragment>
      <Navbar loggedIn={loggedIn} handleDrawerOpen={handleDrawerOpen} />
      <div className={classes.root}>
        <LeftDrawer isOpen={drawerOpen} handleDrawerClose={handleDrawerClose} />
        <main
          className={clsx(classes.content, {
            [classes.contentShift]: drawerOpen,
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

export default Layout;
