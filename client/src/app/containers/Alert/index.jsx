import React from 'react';
import { useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { sliceKey, reducer } from './slice';
import { alertSaga } from './saga';
import { selectAlerts } from './selectors';

const Alert = () => {
  useInjectReducer({ key: sliceKey, reducer });
  useInjectSaga({ key: sliceKey, saga: alertSaga });

  const alerts = useSelector(selectAlerts);

  const alertTypeToStyle = t => {
    switch (t) {
      case 'success':
        return 'alert-success';
      case 'error':
        return 'alert-error';
      case 'warning':
        return 'alert-warning';
      default:
        return 'alert-info';
    }
  };
  return (
    <div className="fixed w-screen px-3 py-2">
      {alerts.map(a => (
        // Issue#61 think about only displaying the last alert
        <div key={a.id} className={`alert ${alertTypeToStyle(a.type)} mb-2`}>
          {a.msg}
        </div>
      ))}
    </div>
  );
};

export default Alert;
