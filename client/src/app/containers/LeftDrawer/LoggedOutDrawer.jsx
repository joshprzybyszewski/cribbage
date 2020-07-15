import React from 'react';
import { useHistory } from 'react-router-dom';
import {
  Link,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  makeStyles,
} from '@material-ui/core';
import CreateIcon from '@material-ui/icons/Create';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';

const useStyles = makeStyles(theme => ({
  link: {
    color: 'inherit',
  },
}));

const LoggedInDrawer = props => {
  const classes = useStyles();
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

export default LoggedInDrawer;
