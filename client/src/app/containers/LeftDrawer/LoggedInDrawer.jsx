import React from 'react';

import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import SearchIcon from '@material-ui/icons/Search';
import CancelIcon from '@material-ui/icons/Cancel';
import HomeIcon from '@material-ui/icons/Home';
import PersonIcon from '@material-ui/icons/Person';
import { authSaga } from 'auth/saga';
import { sliceKey, reducer, actions } from 'auth/slice';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const LoggedInDrawer = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const onClickLogout = () => {
    dispatch(actions.logout(history));
  };
  return (
    <React.Fragment>
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
        <ListItem button onClick={() => history.push('/suggestions')}>
          <ListItemIcon>
            <SearchIcon />
          </ListItemIcon>
          <ListItemText primary='Suggestions' />
        </ListItem>
      </List>
      <Divider />
      <List>
        <ListItem button onClick={onClickLogout}>
          <ListItemIcon>
            <CancelIcon />
          </ListItemIcon>
          <ListItemText primary='Logout' />
        </ListItem>
      </List>
    </React.Fragment>
  );
};

export default LoggedInDrawer;
