import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';

import { authActions } from '../../sagas/actions';

const RegisterForm = ({ register }) => {
  return (
    <Form onFinish={formData => register(formData)}>
      <Form.Item
        name='id'
        label='Username'
        rules={[{ required: true, message: 'Please input your username!' }]}
      >
        <Input placeholder='Username' prefix={<UserOutlined />} />
      </Form.Item>
      <Form.Item
        name='name'
        label='Display Name'
        rules={[{ required: true, message: 'Please input your display name!' }]}
      >
        <Input placeholder='Display name' />
      </Form.Item>
      <Form.Item>
        <Button type='primary' htmlType='submit'>
          Register
        </Button>
      </Form.Item>
    </Form>
  );
};

RegisterForm.propTypes = {
  register: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    register: formData =>
      dispatch(authActions.register(formData.id, formData.name)),
  };
};

export default connect(null, mapDispatchToProps)(RegisterForm);
