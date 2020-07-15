import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Switch, Route } from 'react-router-dom';
import Account from './containers/Account';
import Home from './containers/Home';
import Login from './containers/Login';
import Layout from './containers/Layout';
import NewGameForm from './containers/NewGameForm';
import PrivateRoute from './containers/PrivateRoute';
import Register from './containers/Register';
import '../styles.css';

export const App = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Switch>
          <Route exact path='/' component={Login} />
          <Route exact path='/register' component={Register} />
          <PrivateRoute exact path='/home' component={Home} />
          <PrivateRoute exact path='/newgame' component={NewGameForm} />
          <PrivateRoute exact path='/account' component={Account} />
        </Switch>
      </Layout>
    </BrowserRouter>
  );
};
