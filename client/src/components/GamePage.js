import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const Game = ({ me, gameID }) => {
  return <div>Welcome, {me.name} to game {gameID}! Your ID is {me.id}</div>;
};

Game.propTypes = {
  me: PropTypes.object.isRequired,
  gameID: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  me: state.auth,
  gameID: state.gameID,
});

export default connect(mapStateToProps, null)(Game);
