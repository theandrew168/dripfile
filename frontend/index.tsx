import React, { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import CssBaseline from '@mui/material/CssBaseline';

import { App } from './App';

const app = document.getElementById('app');
if (app) {
	const root = createRoot(app);
	root.render(
		<StrictMode>
			<CssBaseline />
			<App />
		</StrictMode>,
	);
}
