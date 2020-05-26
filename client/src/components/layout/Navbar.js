import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { authActions } from '../../sagas/actions';
import { Button } from 'antd';

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
