import React from 'react';

import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import CancelIcon from '@material-ui/icons/Cancel';
import HomeIcon from '@material-ui/icons/Home';
import PersonIcon from '@material-ui/icons/Person';
import { useHistory } from 'react-router-dom';

import { useAuth } from '../../../auth/useAuth';

const LoggedInDrawer = () => {
    const { logout } = useAuth();
    const history = useHistory();
    return (
        <>
            <List>
                <ListItem button onClick={() => history.push('/home')}>
                    <ListItemIcon>
                        <HomeIcon />
                    </ListItemIcon>
                    <ListItemText primary='Home' />
                </ListItem>
                <ListItem button onClick={() => history.push('/newgame')}>
                    <ListItemIcon>
                        <AddCircleOutlineIcon />
                    </ListItemIcon>
                    <ListItemText primary='New Game' />
                </ListItem>
                <ListItem button onClick={() => history.push('/account')}>
                    <ListItemIcon>
                        <PersonIcon />
                    </ListItemIcon>
                    <ListItemText primary='My Account' />
                </ListItem>
            </List>
            <Divider />
            <List>
                <ListItem button onClick={logout}>
                    <ListItemIcon>
                        <CancelIcon />
                    </ListItemIcon>
                    <ListItemText primary='Logout' />
                </ListItem>
            </List>
        </>
    );
};

export default LoggedInDrawer;
