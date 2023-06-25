import React from "react";

import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { Image } from "./Image";

export function App() {
	const router = createBrowserRouter([
		{
			path: "/",
			element: (
				<div>
					Hello router! <a href="/image">View an image!</a>
				</div>
			),
		},
		{
			path: "/image",
			element: (
				<div>
					Cool image, huh?
					<Image src="/static/logo-black.svg" />
				</div>
			),
		},
	]);

	return <RouterProvider router={router} />;
}
