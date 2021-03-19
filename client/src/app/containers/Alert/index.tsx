import React from 'react';

import { makeStyles, Typography } from '@material-ui/core';
import { green, grey, red, yellow } from '@material-ui/core/colors';
import clsx from 'clsx';
import { useSelector } from 'react-redux';

import { RootState } from '../../../store/store';

const useStyles = makeStyles(theme => ({
    container: {
        position: 'fixed',
        width: '50vw',
        right: 0,
        padding: theme.spacing(0.25, 0.5),
        zIndex: 100,
    },
    alert: {
        // flex items-center h-10 px-4 rounded-lg border-solid border-2;
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 0.5),
        marginBottom: theme.spacing(0.25),
        borderStyle: 'solid',
        borderWidth: 2,
        borderRadius: 4,
    },
    success: {
        backgroundColor: green[200],
        borderColor: green[300],
    },
    error: {
        backgroundColor: red[200],
        borderColor: red[300],
    },
    warning: {
        backgroundColor: yellow[200],
        borderColor: yellow[300],
    },
    info: {
        backgroundColor: grey[400],
        borderColor: grey[500],
    },
}));

const Alert: React.FunctionComponent = () => {
    const alerts = useSelector((state: RootState) => state.alerts);
    const classes = useStyles();
    return (
        <div className={classes.container}>
            {alerts.map(a => (
                // Issue#61 think about only displaying the last alert
                <div
                    key={a.id}
                    className={clsx(classes.alert, classes[a.type])}
                >
                    <Typography variant='h5'>{a.msg}</Typography>
                </div>
            ))}
        </div>
    );
};

export default Alert;
