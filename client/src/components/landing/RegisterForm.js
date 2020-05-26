import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { registerAction } from '../../sagas/auth';

const RegisterForm = ({ registerAsync }) => {
  return (
    <Form onFinish={formData => registerAsync(formData)}>
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
  registerAsync: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    registerAsync: formData =>
      dispatch(registerAction(formData.id, formData.name)),
  };
};

export default connect(null, mapDispatchToProps)(RegisterForm);
