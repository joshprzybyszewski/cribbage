import React from 'react';
import {
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
} from '@material-ui/core';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import CancelIcon from '@material-ui/icons/Cancel';
import HomeIcon from '@material-ui/icons/Home';
import PersonIcon from '@material-ui/icons/Person';

import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';

const LoggedInDrawer = props => {
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
