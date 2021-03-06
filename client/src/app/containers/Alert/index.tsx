import React from 'react';

import { alertSaga } from 'app/containers/Alert/saga';
import { selectAlerts } from 'app/containers/Alert/selectors';
import { sliceKey, reducer } from 'app/containers/Alert/slice';
import { useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

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
    <div className='fixed w-half-screen right-0 px-3 py-2 z-50'>
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
