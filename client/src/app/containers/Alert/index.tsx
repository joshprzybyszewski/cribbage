import React from 'react';

import { useSelector } from 'react-redux';

import { RootState } from '../../../store/store';

const Alert: React.FunctionComponent = () => {
    const alerts = useSelector((state: RootState) => state.alerts);

    return (
        // TODO get rid of tailwind!
        <div className='fixed w-half-screen right-0 px-3 py-2 z-50'>
            {alerts.map(a => (
                // Issue#61 think about only displaying the last alert
                <div key={a.id} className={`alert alert${a.type} mb-2`}>
                    {a.msg}
                </div>
            ))}
        </div>
    );
};

export default Alert;
