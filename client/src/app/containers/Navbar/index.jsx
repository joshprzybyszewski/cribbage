import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { selectLoggedIn } from '../../../auth/selectors';

import { grey } from '@material-ui/core/colors';
import {
  AppBar,
  Button,
  IconButton,
  Link,
  Toolbar,
  Typography,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';

import MenuIcon from '@material-ui/icons/Menu';

const useStyles = makeStyles(theme => ({
  loggedOutLink: {
    color: grey[200],
    marginRight: theme.spacing(2),
  },
  logoutButton: {
    color: grey[200],
    textTransform: 'capitalize',
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
}));

const Navbar = ({ loggedIn }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const onClickLogout = () => {
    dispatch(actions.logout(history));
  };

  const classes = useStyles();

  return (
    <AppBar position='static'>
      <Toolbar>
        {loggedIn && (
          <IconButton
            edge='start'
            className={classes.menuButton}
            color='inherit'
            aria-label='menu'
          >
            <MenuIcon />
          </IconButton>
        )}
        <Typography variant='h6' className={classes.title}>
          Cribbage
        </Typography>
        {loggedIn ? (
          <Button onClick={onClickLogout} className={classes.logoutButton}>
            Logout
          </Button>
        ) : (
          <div>
            <Link href='/' className={classes.loggedOutLink}>
              Login
            </Link>
            <Link href='/register' className={classes.loggedOutLink}>
              Register
            </Link>
          </div>
        )}
      </Toolbar>
    </AppBar>
    // <nav className='h-12 px-4 bg-blue-900 flex justify-between items-center text-gray-400'>
    //   <Link
    //     to={loggedIn ? '/home' : '/'}
    //     className='uppercase text-xl tracking-wider hover:text-white'
    //   >
    //     Cribbage
    //   </Link>
    //   {!loggedIn ? (
    //     <div className='flex'>
    //       <Link to='/' className='px-2 hover:text-white'>
    //         Login
    //       </Link>
    //       <Link to='/register' className='px-2 hover:text-white'>
    //         Register
    //       </Link>
    //     </div>
    //   ) : (
    //     <button
    //       onClick={onClickLogout}
    //       className='focus:outline-none hover:text-white'
    //     >
    //       Logout
    //     </button>
    //   )}
    // </nav>
  );
};

export default Navbar;
