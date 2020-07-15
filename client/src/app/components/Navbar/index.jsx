import React from 'react';
import { grey } from '@material-ui/core/colors';
import { AppBar, IconButton, Toolbar, Typography } from '@material-ui/core';
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

const Navbar = ({ handleDrawerOpen }) => {
  const classes = useStyles();

  return (
    <AppBar position='static'>
      <Toolbar>
        <IconButton
          edge='start'
          className={classes.menuButton}
          color='inherit'
          aria-label='menu'
          onClick={handleDrawerOpen}
        >
          <MenuIcon />
        </IconButton>
        <Typography variant='h6' className={classes.title}>
          Cribbage
        </Typography>
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;
