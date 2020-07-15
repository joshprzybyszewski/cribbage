import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import { useInjectReducer } from 'redux-injectors';
import clsx from 'clsx';
import { Divider, Drawer, IconButton } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import { selectLoggedIn } from '../../../auth/selectors';
import { sliceKey, reducer } from '../../../auth/slice';
import Alert from '../Alert';
import Navbar from '../../components/Navbar';
import LoggedInDrawer from '../LeftDrawer/LoggedInDrawer';
import LoggedOutDrawer from '../LeftDrawer/LoggedOutDrawer';

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
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
  const handleDrawerOpen = () => {
    setDrawerOpen(true);
  };
  const handleDrawerClose = () => {
    setDrawerOpen(false);
  };

  const classes = useStyles();

  return (
    <React.Fragment>
      <Navbar loggedIn={loggedIn} handleDrawerOpen={handleDrawerOpen} />
      <div className={classes.root}>
        <Drawer
          className={classes.drawer}
          variant='persistent'
          anchor='left'
          open={drawerOpen}
          classes={{
            paper: classes.drawerPaper,
          }}
        >
          <div className={classes.drawerHeader}>
            <IconButton onClick={handleDrawerClose}>
              <ChevronLeftIcon />
            </IconButton>
          </div>
          <Divider />
          {loggedIn ? <LoggedInDrawer /> : <LoggedOutDrawer />}
        </Drawer>
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

export default Layout;
