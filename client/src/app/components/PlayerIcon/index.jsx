import React from 'react';

import blue from '@material-ui/core/colors/blue';
import green from '@material-ui/core/colors/green';
import grey from '@material-ui/core/colors/grey';
import red from '@material-ui/core/colors/red';
import PersonPinCircleIcon from '@material-ui/icons/PersonPinCircle';
import PropTypes from 'prop-types';

const colors = {
  red: red[800],
  blue: blue[800],
  green: green[800],
};

const colorStringToColor = str => {
  if (Object.keys(colors).includes(str)) {
    return colors[str];
  }
  return grey[400];
};

const PlayerIcon = ({ color }) => (
  <PersonPinCircleIcon style={{ color: colorStringToColor(color) }} />
);

PlayerIcon.propTypes = {
  color: PropTypes.string.isRequired,
};

export default PlayerIcon;
