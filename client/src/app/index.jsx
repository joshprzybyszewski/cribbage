import React from 'react';
import { ConnectedRouter } from 'connected-react-router';
import { Switch, Route } from 'react-router-dom';
import { history } from '../store/reducers';
import Home from './containers/Home';
import Register from './containers/Register';
import PrivateRoute from './containers/PrivateRoute';
import Login from './containers/Login';

export const App = () => {
  return (
    <ConnectedRouter history={history}>
      <div className='relative bg-gray-200 h-screen'>
        <Navbar />
        <Alert />
        <Switch>
          <Route exact path='/' component={Register} />
          <Route exact path='/login' component={Login} />
          <PrivateRoute exact path='/home' component={Home} />
        </Switch>
      </div>
    </ConnectedRouter>
  );
};
