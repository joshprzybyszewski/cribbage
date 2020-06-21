import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button } from 'antd';
import { gameActions } from '../sagas/actions';

const Home = ({ currentPlayer, allgames, goToGame }) => {
  return <div>
    Welcome, {currentPlayer.name}! Your games are:
    <div>
      {allgames ? allgames.map((gID, desc) => (
        <Button key={gID} onClick={() => goToGame({gID})}>{desc}</Button>
      )) : 'nothing'}
      {/* <Button onClick={() => goToGame(1863140844)}>Game 1863140844</Button> <br></br>
      <Button onClick={() => goToGame(456)}>Game 456</Button> <br></br> */}
    </div>
  </div>;
};

Home.propTypes = {
  currentPlayer: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  currentPlayer: state.auth,
  allgames: state.allgames,
});

const mapDispatchToProps = dispatch => {
  return {
    goToGame: gID => dispatch(gameActions.viewGame(gID)),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Home);
