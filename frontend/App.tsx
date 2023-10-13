import React from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import Image from "./Image";
import ErrorPage from "./ErrorPage";
import LocationList from "./LocationList";
import LocationRead from "./LocationRead";

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
					<Image src="/static/img/logo-black.svg" />
				</div>
			),
		},
		{
			path: "/location",
			element: <LocationList />,
		},
		{
			path: "/location/:id",
			element: <LocationRead />,
		},
	]);

	return <RouterProvider router={router} />;
}
