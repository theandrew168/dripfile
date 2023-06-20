import React from 'react';

import { NavigationBar } from './NavigationBar';
import { Image } from './Image';
import { CoolButton } from './CoolButton';

export function App() {
	return (
		<>
			<NavigationBar />
			<CoolButton message='Hello World' />
			<Image src="/static/logo-black.svg" />
		</>
	);
}
