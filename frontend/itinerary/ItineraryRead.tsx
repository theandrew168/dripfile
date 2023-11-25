import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

import type { ItineraryReadResponse } from "../types";

export default function ItineraryRead() {
	const { id } = useParams();

	const { isPending, isError, error, data } = useQuery({
		queryKey: ["itineraries", id],
		queryFn: async () => {
			const response = await fetch(`/api/v1/itineraries/${id}`);
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: ItineraryReadResponse = await response.json();
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

	const itinerary = data.itinerary;
	return (
		<div>
			<p>ID: {itinerary.id}</p>
			<p>From: {itinerary.fromLocationID}</p>
			<p>To: {itinerary.toLocationID}</p>
			<p>Pattern: {itinerary.pattern}</p>
		</div>
	);
}
