import React from "react";
import { Link } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { isErrorResponse, type LocationListResponse } from "../types";
import Alert from "../Alert";

export default function ItineraryCreate() {
	const queryClient = useQueryClient();
	const { mutate, isPending, isError, error, isSuccess } = useMutation({
		mutationFn: async (form: FormData) => {
			const payload = {
				...Object.fromEntries(form),
			};
			const response = await fetch("/api/v1/itineraries", {
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
			queryClient.invalidateQueries({ queryKey: ["itineraries"] });
		},
	});

	const locationsQuery = useQuery({
		queryKey: ["locations"],
		queryFn: async () => {
			const response = await fetch("/api/v1/locations");
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: LocationListResponse = await response.json();
			return data;
		},
	});

	// TODO: build a generic loading component
	if (locationsQuery.isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (locationsQuery.isError) {
		return <div>Error: {locationsQuery.error.message}</div>;
	}

	const locations = locationsQuery.data.locations;

	// https://tailwindui.com/components/application-ui/forms/form-layouts#component-dcf2bee8aa4fbef0d4623df5b9718da8
	return (
		<>
			{isError && <Alert type="failure" message={error.message} />}
			{isSuccess && <Alert type="success" message="Itinerary created!" />}
			<div className="space-y-10 divide-y divide-gray-900/10">
				<div className="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
					<div className="px-4 sm:px-0">
						<h2 className="text-base font-semibold leading-7 text-gray-900">Itinerary</h2>
						<p className="mt-1 text-sm leading-6 text-gray-600">A plan for transferring files between locations.</p>
					</div>

					<form
						onSubmit={(event) => {
							event.preventDefault();
							mutate(new FormData(event.currentTarget));
						}}
						className="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2"
					>
						<div className="px-4 py-6 sm:p-8">
							<div className="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
								<div className="sm:col-span-4">
									<label htmlFor="pattern" className="block text-sm font-medium leading-6 text-gray-900">
										Pattern
									</label>
									<div className="mt-2">
										<div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
											<input
												type="text"
												name="pattern"
												id="pattern"
												className="block flex-1 border-0 bg-transparent py-1.5 pl-2 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
											/>
										</div>
									</div>
								</div>
								<div className="sm:col-span-4">
									<label htmlFor="fromLocationID" className="block text-sm font-medium leading-6 text-gray-900">
										From
									</label>
									<div className="mt-2">
										<select
											id="fromLocationID"
											name="fromLocationID"
											className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:max-w-xs sm:text-sm sm:leading-6"
										>
											{locations.map((location) => (
												<option key={location.id} value={location.id}>
													{location.id}
												</option>
											))}
										</select>
									</div>
								</div>
								<div className="sm:col-span-4">
									<label htmlFor="toLocationID" className="block text-sm font-medium leading-6 text-gray-900">
										To
									</label>
									<div className="mt-2">
										<select
											id="toLocationID"
											name="toLocationID"
											className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:max-w-xs sm:text-sm sm:leading-6"
										>
											{locations.map((location) => (
												<option key={location.id} value={location.id}>
													{location.id}
												</option>
											))}
										</select>
									</div>
								</div>
							</div>
						</div>
						<div className="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8">
							<Link to="/itineraries">
								<button type="button" className="text-sm font-semibold leading-6 text-gray-900">
									Cancel
								</button>
							</Link>
							<button
								type="submit"
								disabled={isPending}
								className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
							>
								{isPending ? "Creating..." : "Create"}
							</button>
						</div>
					</form>
				</div>
			</div>
		</>
	);
}
