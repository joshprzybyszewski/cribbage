import React, { Fragment } from 'react';
import { Link } from 'react-router-dom';

const Navbar = props => (
  <div>
    <Link to='/home' className='link'>
      <h1 style={{ color: '#fff' }}>CRIBBAGE</h1>
    </Link>
  </div>
);

export default Navbar;
