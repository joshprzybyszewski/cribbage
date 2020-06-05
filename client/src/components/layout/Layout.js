import React from 'react';

import Alert from './Alert';
import Navbar from './Navbar';

const Layout = props => {
  return (
    <div className='bg-gray-200 h-screen'>
      <Navbar />
      <Alert />
      {props.children}
    </div>
  );
};

export default Layout;
