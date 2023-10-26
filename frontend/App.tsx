import React from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import ErrorPage from "./ErrorPage";
import ItineraryPage from "./itinerary/ItineraryPage";
import ItineraryCreate from "./itinerary/ItineraryCreate";
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
					path: "/locations",
					element: <LocationPage />,
				},
				{
					path: "/locations/create",
					element: <LocationCreate />,
				},
				{
					path: "/locations/create/in-memory",
					element: <LocationCreateInMemory />,
				},
				{
					path: "/locations/create/s3",
					element: <LocationCreateS3 />,
				},
				{
					path: "/locations/:id",
					element: <LocationRead />,
				},
				{
					path: "/itineraries",
					element: <ItineraryPage />,
				},
				{
					path: "/itineraries/create",
					element: <ItineraryCreate />,
				},
			],
		},
	]);

	return <RouterProvider router={router} />;
}
