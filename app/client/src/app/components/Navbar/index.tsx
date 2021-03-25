import React from 'react';

import AppBar from '@material-ui/core/AppBar';
import IconButton from '@material-ui/core/IconButton';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import MenuIcon from '@material-ui/icons/Menu';

const useStyles = makeStyles(theme => ({
    menuButton: {
        marginRight: theme.spacing(2),
    },
    title: {
        flexGrow: 1,
    },
}));

interface Props {
    handleDrawerOpen: () => void;
}

const Navbar: React.FunctionComponent<Props> = ({ handleDrawerOpen }) => {
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
