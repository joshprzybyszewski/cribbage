import React from 'react';

import Account from 'app/containers/Account';
import Game from 'app/containers/Game';
import Home from 'app/containers/Home';
import Layout from 'app/containers/Layout';
import Login from 'app/containers/Login';
import NewGameForm from 'app/containers/NewGameForm';
import PrivateRoute from 'app/containers/PrivateRoute';
import Register from 'app/containers/Register';
import { Switch, Route } from 'react-router-dom';
import { BrowserRouter } from 'react-router-dom';
import 'styles.css';

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
          <PrivateRoute exact path='/game' component={Game} />
        </Switch>
      </Layout>
    </BrowserRouter>
  );
};
