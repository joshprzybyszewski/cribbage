import React from 'react';

import { Switch, Route, BrowserRouter } from 'react-router-dom';

import Account from './containers/Account';
import Game from './containers/Game';
import Home from './containers/Home';
import Layout from './containers/Layout';
import Login from './containers/Login';
import NewGameForm from './containers/NewGameForm';
import PrivateRoute from './containers/PrivateRoute';
import Register from './containers/Register';
import Suggestions from './containers/Suggestions';

export const App = () => {
    return (
        <BrowserRouter>
            <Layout>
                <Switch>
                    <Route exact path='/' component={Login} />
                    <Route exact path='/register' component={Register} />
                    <PrivateRoute exact path='/home' component={Home} />
                    <PrivateRoute
                        exact
                        path='/newgame'
                        component={NewGameForm}
                    />
                    <PrivateRoute exact path='/account' component={Account} />
                    <PrivateRoute exact path='/game' component={Game} />
                    <PrivateRoute exact path='/suggestions' component={Suggestions} />
                </Switch>
            </Layout>
        </BrowserRouter>
    );
};
