import React from 'react';
import { useSelector } from 'react-redux';
import { useInjectReducer } from 'redux-injectors';
import { selectLoggedIn } from '../auth/selectors';
import { sliceKey, reducer } from '../auth/slice';

import { BrowserRouter } from 'react-router-dom';
import { Switch, Route } from 'react-router-dom';
import Account from './containers/Account';
import Alert from './containers/Alert';
import Home from './containers/Home';
import Login from './containers/Login';
import Navbar from './containers/Navbar';
import NewGameForm from './containers/NewGameForm';
import PrivateRoute from './containers/PrivateRoute';
import Register from './containers/Register';
import Sidebar from './components/Sidebar';
import '../styles.css';

export const App = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  const loggedIn = useSelector(selectLoggedIn);

  return (
    <BrowserRouter>
      <Navbar />
      <div className='flex flex-row'>
        {loggedIn && <Sidebar />}
        <div className='w-full h-screen p-4 bg-gray-200'>
          <Alert />
          <Switch>
            <Route exact path='/' component={Login} />
            <Route exact path='/register' component={Register} />
            <PrivateRoute exact path='/home' component={Home} />
            <PrivateRoute exact path='/newgame' component={NewGameForm} />
            <PrivateRoute exact path='/account' component={Account} />
          </Switch>
        </div>
      </div>
    </BrowserRouter>
  );
};
