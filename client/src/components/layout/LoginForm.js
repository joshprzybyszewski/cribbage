import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { LOGIN_ASYNC } from '../../sagas/types';

const LoginForm = ({ loginAsync }) => {
  const onFinish = formData => {
    loginAsync(formData);
  };
  return (
    <Form onFinish={onFinish}>
      <Form.Item
        name='username'
        label='Username'
        rules={[{ required: true, message: 'Please input your username!' }]}
      >
        <Input placeholder='Username' prefix={<UserOutlined />} />
      </Form.Item>
      <Form.Item>
        <Button type='primary' htmlType='submit'>
          Login
        </Button>
      </Form.Item>
    </Form>
  );
};

LoginForm.propTypes = {
  loginAsync: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    loginAsync: formData => dispatch({ type: LOGIN_ASYNC, payload: formData }),
  };
};

export default connect(null, mapDispatchToProps)(LoginForm);
