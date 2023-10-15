import React from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import ErrorPage from "./ErrorPage";
import LocationPage from "./LocationPage";
import LocationRead from "./LocationRead";
import NavBar from "./NavBar";

export default function App() {
	const router = createBrowserRouter([
		{
			path: "/",
			element: <NavBar />,
			errorElement: <ErrorPage />,
			children: [
				{
					path: "/locations",
					element: <LocationPage />,
				},
				{
					path: "/locations/:id",
					element: <LocationRead />,
				},
			],
		},
	]);

	return <RouterProvider router={router} />;
}
