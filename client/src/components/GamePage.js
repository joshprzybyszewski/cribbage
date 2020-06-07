import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const Game = ({ me, game }) => {
  return <div>
    Welcome, {me.name} to game {game.gameID}!
    Your ID is {me.id}.
    Players are {game.players.toString()}.
    Phase is {game.phase}.
    </div>;
};

Game.propTypes = {
  me: PropTypes.object.isRequired,
  game: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  me: state.auth,
  game: state.game,
});

export default connect(mapStateToProps, null)(Game);
