import React from 'react';

import AppBar from '@material-ui/core/AppBar';
import grey from '@material-ui/core/colors/grey';
import IconButton from '@material-ui/core/IconButton';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import PropTypes from 'prop-types';

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

Navbar.propTypes = {
  handleDrawerOpen: PropTypes.func.isRequired,
};

export default Navbar;
