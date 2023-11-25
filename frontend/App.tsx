import React from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import ErrorPage from "./ErrorPage";
import ItineraryCreate from "./itinerary/ItineraryCreate";
import ItineraryPage from "./itinerary/ItineraryPage";
import ItineraryRead from "./itinerary/ItineraryRead";
import LocationCreate from "./location/LocationCreate";
import LocationCreateInMemory from "./location/LocationCreateInMemory";
import LocationCreateS3 from "./location/LocationCreateS3";
import LocationPage from "./location/LocationPage";
import LocationRead from "./location/LocationRead";
import NavBar from "./NavBar";

export default function App() {
	const router = createBrowserRouter([
		{
			path: "/",
			element: <NavBar />,
			errorElement: <ErrorPage />,
			children: [
				{
					path: "/location",
					element: <LocationPage />,
				},
				{
					path: "/location/create",
					element: <LocationCreate />,
				},
				{
					path: "/location/create/in-memory",
					element: <LocationCreateInMemory />,
				},
				{
					path: "/location/create/s3",
					element: <LocationCreateS3 />,
				},
				{
					path: "/location/:id",
					element: <LocationRead />,
				},
				{
					path: "/itinerary",
					element: <ItineraryPage />,
				},
				{
					path: "/itinerary/create",
					element: <ItineraryCreate />,
				},
				{
					path: "/itinerary/:id",
					element: <ItineraryRead />,
				},
			],
		},
	]);

	return <RouterProvider router={router} />;
}
