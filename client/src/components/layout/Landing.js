import React from 'react';
import { Button, Space } from 'antd';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { login, register } from '../../actions/auth';

const Landing = props => {
  return (
    <section className='landing'>
      <div className='dark-overlay'>
        <div className='landing-inner'>
          <h1>Welcome to Cribbage!</h1>
          <p>Login or register to play cribbage against your friends online</p>
          <Space>
            <Button size='large' type='primary'>
              Login
            </Button>
            <Button size='large'>Register</Button>
          </Space>
        </div>
      </div>
    </section>
  );
};

Landing.propTypes = {
  login: PropTypes.func.isRequired,
  register: PropTypes.func.isRequired,
};

export default connect(null, { login, register })(Landing);
