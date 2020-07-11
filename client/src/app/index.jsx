import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Switch, Route } from 'react-router-dom';
import Alert from './containers/Alert';
import Game from './containers/Game';
import Home from './containers/Home';
import Login from './containers/Login';
import Navbar from './containers/Navbar';
import PrivateRoute from './containers/PrivateRoute';
import Register from './containers/Register';
import '../styles.css';

export const App = () => {
  return (
    <BrowserRouter>
      <div className='relative bg-gray-200 h-screen'>
        <Navbar />
        <Alert />
        <Switch>
          <Route exact path='/' component={Login} />
          <Route exact path='/register' component={Register} />
          <PrivateRoute exact path='/home' component={Home} />
          <PrivateRoute exact path='/game' component={Game} />
        </Switch>
      </div>
    </BrowserRouter>
  );
};
