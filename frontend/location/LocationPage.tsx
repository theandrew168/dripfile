import React from "react";
import { useQuery } from "@tanstack/react-query";

import LocationEmpty from "./LocationEmpty";
import LocationList from "./LocationList";
import { listLocations } from "../fetch";

export default function LocationPage() {
	const {
		isPending,
		isError,
		error,
		data: locations,
	} = useQuery({
		queryKey: ["location"],
		queryFn: async () => listLocations(),
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const hasLocations = locations.length > 0;
	return hasLocations ? <LocationList locations={locations} /> : <LocationEmpty />;
}
