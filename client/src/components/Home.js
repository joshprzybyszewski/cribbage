import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const Home = ({ currentPlayer }) => {
  return <div>Welcome, {currentPlayer.name}!</div>;
};

Home.propTypes = {
  currentPlayer: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  currentPlayer: state.auth,
});

export default connect(mapStateToProps, null)(Home);
