import React from "react";
import { useQuery } from "@tanstack/react-query";

import type { LocationListResponse } from "../types";
import LocationEmpty from "./LocationEmpty";
import LocationList from "./LocationList";

export default function LocationPage() {
	const { isPending, isError, error, data } = useQuery({
		queryKey: ["location"],
		queryFn: async () => {
			const response = await fetch("/api/v1/location");
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: LocationListResponse = await response.json();
			return data;
		},
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const locations = data.locations;
	const hasLocations = locations.length > 0;
	return hasLocations ? <LocationList locations={locations} /> : <LocationEmpty />;
}
