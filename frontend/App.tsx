import React from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import Image from "./Image";
import ErrorPage from "./ErrorPage";

export default function App() {
	const router = createBrowserRouter([
		{
			path: "/",
			element: (
				<div>
					Hello router! <a href="/image">View an image!</a>
				</div>
			),
			errorElement: <ErrorPage />,
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
