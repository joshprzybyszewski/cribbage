import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

const Home = ({ currentPlayer }) => {
  return <div>
    Welcome, {currentPlayer.name}! Your games are:
    <div>
      <Link to='/game/123'>Game 123</Link>
    </div>
  </div>;
};

Home.propTypes = {
  currentPlayer: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  currentPlayer: state.auth,
});

export default connect(mapStateToProps, null)(Home);
