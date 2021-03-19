import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import makeStyles from '@material-ui/core/styles/makeStyles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { useHistory } from 'react-router-dom';

import { useAuth } from '../../../auth/useAuth';
import { useGame } from '../Game/useGame';

const useStyles = makeStyles(theme => ({
    title: {
        fontSize: '2rem',
    },
    paper: {
        marginTop: theme.spacing(8),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
    },
    form: {
        marginTop: theme.spacing(1),
    },
    submit: {
        margin: theme.spacing(3, 0, 2),
    },
}));

interface FormData {
    [key: string]: string;
}

const NewGameForm: React.FunctionComponent = () => {
    const history = useHistory();
    const { currentUser } = useAuth();
    const { createGame } = useGame();

    // event handlers
    const [formData, setFormData] = useState<FormData>({
        id1: '',
        id2: '',
        teammateID: '',
    });
    const onFormDataChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };
    const onSubmitLoginForm = async (e: React.FormEvent) => {
        e.preventDefault();

        const playerIDs = [
            currentUser.id,
            ...Object.keys(formData)
                .map(k => formData[k])
                .filter(id => id.length > 0),
        ];
        await createGame(playerIDs);
        history.push('/game');
    };

    const classes = useStyles();

    return (
        <Container component='main' maxWidth='sm'>
            <CssBaseline />
            <div className={classes.paper}>
                <Typography component='h1' className={classes.title}>
                    Start Game
                </Typography>
                <form className={classes.form} onSubmit={onSubmitLoginForm}>
                    <TextField
                        variant='outlined'
                        margin='normal'
                        required
                        fullWidth
                        label='Opponent 1'
                        name='id1'
                        autoFocus
                        onChange={onFormDataChange}
                    />
                    <TextField
                        disabled
                        variant='outlined'
                        margin='normal'
                        fullWidth
                        label='Opponent 2'
                        name='id2'
                        autoFocus
                        onChange={onFormDataChange}
                    />
                    <TextField
                        disabled
                        variant='outlined'
                        margin='normal'
                        fullWidth
                        label='Teammate'
                        name='teammateID'
                        autoFocus
                        onChange={onFormDataChange}
                    />
                    <Button
                        type='submit'
                        fullWidth
                        variant='contained'
                        color='primary'
                        className={classes.submit}
                        disabled={formData.id1.length === 0}
                    >
                        Challenge
                    </Button>
                </form>
            </div>
        </Container>
    );
};

export default NewGameForm;
