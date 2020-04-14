import React from 'react';
import { Button, Space } from 'antd';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { LOGIN_ASYNC, REGISTER_ASYNC } from '../../sagas/types';

const Landing = ({ loginAsync, registerAsync }) => {
  return (
    <section className="landing">
      <div className="dark-overlay">
        <div className="landing-inner">
          <h1>Welcome to Cribbage!</h1>
          <p>Login or register to play cribbage against your friends online</p>
          <Space>
            <Button onClick={e => loginAsync()} size="large" type="primary">
              Login
            </Button>
            <Button onClick={e => registerAsync()} size="large">
              Register
            </Button>
          </Space>
        </div>
      </div>
    </section>
  );
};

Landing.propTypes = {
  loginAsync: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    loginAsync: () => dispatch({ type: LOGIN_ASYNC, payload: `user!` }),
    registerAsync: () => dispatch({ type: REGISTER_ASYNC, payload: `user!` }),
  };
};

export default connect(null, mapDispatchToProps)(Landing);
