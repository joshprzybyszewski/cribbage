import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Alert as AntAlert } from 'antd';

const Alert = ({ alerts }) => {
  return (
    <div>
      {alerts.map(a => (
        <AntAlert key={a.id} message={a.msg} type={a.type} banner />
      ))}
    </div>
  );
};

Alert.propTypes = {
  alerts: PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
  alerts: state.alert,
});

export default connect(mapStateToProps, null)(Alert);
