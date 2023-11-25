import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

import type { LocationReadResponse } from "../types";

export default function LocationRead() {
	const { id } = useParams();

	const { isPending, isError, error, data } = useQuery({
		queryKey: ["location", id],
		queryFn: async () => {
			const response = await fetch(`/api/v1/location/${id}`);
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: LocationReadResponse = await response.json();
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

	const location = data.location;
	return (
		<div>
			<p>ID: {location.id}</p>
			<p>Kind: {location.kind}</p>
		</div>
	);
}
