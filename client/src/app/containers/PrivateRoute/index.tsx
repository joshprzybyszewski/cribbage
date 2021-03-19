import React from 'react';

import { Redirect, Route, RouteProps } from 'react-router-dom';

import { useAuth } from '../../../auth/useAuth';

interface Props extends RouteProps {}

const PrivateRoute: React.FunctionComponent<Props> = props => {
    const { isLoggedIn } = useAuth();
    if (!isLoggedIn) {
        return <Redirect to='/' />;
    }
    return <Route {...props} />;
};

export default PrivateRoute;
