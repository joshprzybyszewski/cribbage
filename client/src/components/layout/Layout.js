import React from 'react';
import { Layout as AntLayout } from 'antd';
import Alert from './Alert';
import Navbar from './Navbar';

const { Content, Footer, Header } = AntLayout;

const Layout = props => {
  return (
    <AntLayout>
      <Header className='header'>
        <Navbar />
      </Header>
      <Alert />
      <Content className='content'>{props.children}</Content>
      <Footer className='footer'>Footer</Footer>
    </AntLayout>
  );
};

export default Layout;
