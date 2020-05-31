import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { Switch, Route } from 'react-router-dom';

import configureStore, { history } from './store';
import Home from './components/Home';
import Game from './components/GamePage';
import Landing from './components/landing/Landing';
import Layout from './components/layout/Layout';
import PrivateRoute from './components/layout/PrivateRoute';
import './App.css';

const store = configureStore({});

function App() {
  return (
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <Layout>
          <Switch>
            <Route exact path='/' component={Landing} />
            <PrivateRoute exact path='/home' component={Home} />
            <PrivateRoute path='/game' component={Game} />
          </Switch>
        </Layout>
      </ConnectedRouter>
    </Provider>
  );
}

export default App;
