import * as React from 'react';

import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import { App } from './app';
import { configureAppStore } from './store/configureStore';

const store = configureAppStore();
const MOUNT_NODE = document.getElementById('root');

const ConnectedApp = ({ Component }) => (
  <Provider store={store}>
    <React.StrictMode>
      <Component />
    </React.StrictMode>
  </Provider>
);

const render = Component => {
  ReactDOM.render(<ConnectedApp Component={Component} />, MOUNT_NODE);
};

if (module.hot) {
  module.hot.accept(['./app'], () => {
    ReactDOM.unmountComponentAtNode(MOUNT_NODE);
    const App = require('./app').App;
    render(App);
  });
}

render(App);
