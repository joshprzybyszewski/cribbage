import React from 'react';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import grey from '@material-ui/core/colors/grey';
import IconButton from '@material-ui/core/IconButton';
import Link from '@material-ui/core/Link';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import { authSaga } from 'auth/saga';
import { sliceKey, reducer, actions } from 'auth/slice';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

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

const Navbar = ({ loggedIn, handleDrawerOpen }) => {
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
                {loggedIn ? (
                    <Button
                        onClick={onClickLogout}
                        className={classes.logoutButton}
                    >
                        Logout
                    </Button>
                ) : (
                    <div>
                        <Link href='/' className={classes.loggedOutLink}>
                            Login
                        </Link>
                        <Link
                            href='/register'
                            className={classes.loggedOutLink}
                        >
                            Register
                        </Link>
                    </div>
                )}
            </Toolbar>
        </AppBar>
    );
};

Navbar.propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    handleDrawerOpen: PropTypes.func.isRequired,
};

export default Navbar;
