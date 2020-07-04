import React from 'react';
import { useSelector } from 'react-redux';
import { Redirect, Route } from 'react-router-dom';
import { selectLoggedIn } from '../../../auth/selectors';

const PrivateRoute = ({ component: Component, ...rest }) => {
  const loggedIn = useSelector(selectLoggedIn);
  return (
    <Route
      {...rest}
      render={props =>
        loggedIn ? <Component {...props} /> : <Redirect to='/' />
      }
    />
  );
};

export default PrivateRoute;
