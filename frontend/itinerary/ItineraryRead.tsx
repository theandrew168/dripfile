import React from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";

import { isErrorResponse } from "../types";
import { readItinerary } from "../fetch";

export default function ItineraryRead() {
	const { id } = useParams();
	if (!id) {
		return null;
	}

	const {
		isPending,
		isError,
		error,
		data: itinerary,
	} = useQuery({
		queryKey: ["itinerary", id],
		queryFn: async () => readItinerary(id),
	});

	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const mutation = useMutation({
		mutationFn: async (form: FormData) => {
			const payload = {
				...Object.fromEntries(form),
			};
			const response = await fetch("/api/v1/transfer", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(payload),
			});
			if (!response.ok) {
				const error = await response.json();
				if (isErrorResponse(error)) {
					throw new Error(JSON.stringify(error.error));
				} else {
					throw new Error("Network response was not OK");
				}
			}

			return response.json();
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["transfer"] });
			navigate("/transfer");
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

	return (
		<div>
			<p>ID: {itinerary.id}</p>
			<p>From: {itinerary.fromLocationID}</p>
			<p>To: {itinerary.toLocationID}</p>
			<p>Pattern: {itinerary.pattern}</p>
			<p>CreatedAt: {itinerary.createdAt.toString()}</p>
			<p>UpdatedAt: {itinerary.updatedAt.toString()}</p>
			<form
				onSubmit={(event) => {
					event.preventDefault();
					mutation.mutate(new FormData(event.currentTarget));
				}}
			>
				<input type="hidden" name="itineraryID" id="itineraryID" value={itinerary.id} />
				<button
					type="submit"
					className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
				>
					Run Now
				</button>
			</form>
		</div>
	);
}
