import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { Button } from 'antd';

import { authActions } from '../../sagas/actions';

const Navbar = ({ currentUser, logout }) => {
  return (
    <div>
      <Link to='/home'>Home</Link>
      {currentUser.id && <Button onClick={logout}>Logout</Button>}
    </div>
  );
};

Navbar.propTypes = {
  currentPlayer: PropTypes.object,
};

const mapStateToProps = state => ({
  currentUser: state.auth,
});

const mapDispatchToProps = dispatch => {
  return {
    logout: () => dispatch(authActions.logout()),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Navbar);
