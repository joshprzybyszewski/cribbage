import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const Home = ({ player }) => {
  return (
    <div>
      <h1 className='center xl'>Welcome, {player.n}!</h1>
    </div>
  );
};

Home.propTypes = {
  player: PropTypes.object.isRequired,
};

const mapStateToProps = state => ({
  player: state.auth.player,
});

export default connect(mapStateToProps, null)(Home);
