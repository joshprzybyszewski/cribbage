import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import { authActions } from '../../sagas/actions';

const Navbar = ({ currentUser, logout }) => {
  return (
    <nav className='h-12 px-4 bg-blue-900 flex justify-between items-center text-gray-400'>
      <Link
        to='/'
        className='uppercase text-xl tracking-wider hover:text-white'
      >
        Cribbage
      </Link>
      <div className='flex'>
        <Link to='/login' className='px-2 hover:text-white'>
          Login
        </Link>
        <Link to='/register' className='px-2 hover:text-white'>
          Register
        </Link>
      </div>
    </nav>
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
