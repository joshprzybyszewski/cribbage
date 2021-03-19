import { blue, green, grey, red } from '@material-ui/core/colors';

export const colorToHue = (color: string) => {
    if (color.includes('red')) {
        return red[800];
    }
    if (color.includes('green')) {
        return green[800];
    }
    if (color.includes('blue')) {
        return blue[800];
    }
    return grey[400];
};
