import React, { StrictMode } from 'react'
import { createRoot } from 'react-dom/client';

import { NavigationBar } from './NavigationBar';

const domNode = document.getElementById('navigation');
if (domNode) {
	const root = createRoot(domNode);
	root.render(
		<StrictMode>
			<NavigationBar />			
		</StrictMode>
	);
}
