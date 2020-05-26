import React from 'react';
import { Layout as AntLayout } from 'antd';
import Alert from './Alert';

const { Content, Footer, Header } = AntLayout;

const Layout = props => {
  return (
    <AntLayout>
      <Header className='header'>Header</Header>
      <Alert />
      <Content className='content'>{props.children}</Content>
      <Footer className='footer'>Footer</Footer>
    </AntLayout>
  );
};

export default Layout;
