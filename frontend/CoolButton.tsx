import React from 'react';
import Button from '@mui/material/Button';

type Props = {
	message: string;
};

export function CoolButton({ message }: Props) {
	return (
		<Button variant="contained">{message}</Button>
	);
}
