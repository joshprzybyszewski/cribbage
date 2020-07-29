import React from 'react';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import CreateIcon from '@material-ui/icons/Create';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import { useHistory } from 'react-router-dom';

const LoggedOutDrawer = () => {
  const history = useHistory();

  return (
    <List>
      <ListItem button onClick={() => history.push('/')}>
        <ListItemIcon>
          <ExitToAppIcon />
        </ListItemIcon>
        <ListItemText primary='Login' />
      </ListItem>
      <ListItem button onClick={() => history.push('/register')}>
        <ListItemIcon>
          <CreateIcon />
        </ListItemIcon>
        <ListItemText primary='Register' />
      </ListItem>
    </List>
  );
};

export default LoggedOutDrawer;
