import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import Alert from "../Alert";
import { createInMemoryLocation } from "../fetch";
import type { NewInMemoryLocation } from "../types";

export default function LocationCreateInMemory() {
	const [capacity, setCapacity] = useState("");

	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const { mutate, isPending, isError, error } = useMutation({
		mutationFn: (location: NewInMemoryLocation) => createInMemoryLocation(location),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["location"] });
			navigate("/location");
		},
	});

	// https://tailwindui.com/components/application-ui/forms/form-layouts#component-dcf2bee8aa4fbef0d4623df5b9718da8
	return (
		<>
			{isError && <Alert type="failure" message={error.message} />}
			<div className="space-y-10 divide-y divide-gray-900/10">
				<div className="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
					<div className="px-4 sm:px-0">
						<h2 className="text-base font-semibold leading-7 text-gray-900">In-Memory Location</h2>
						<p className="mt-1 text-sm leading-6 text-gray-600">An in-memory location for testing Dripfile.</p>
					</div>

					<form
						onSubmit={(event) => {
							event.preventDefault();
							event.stopPropagation();
							mutate({ kind: "memory", capacity });
						}}
						className="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2"
					>
						<div className="px-4 py-6 sm:p-8">
							<div className="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
								<div className="sm:col-span-4">
									<label htmlFor="capcity" className="block text-sm font-medium leading-6 text-gray-900">
										Capacity
									</label>
									<div className="mt-2">
										<div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
											<input
												type="text"
												id="capacity"
												name="capacity"
												value={capacity}
												onChange={(event) => setCapacity(event.target.value)}
												className="block flex-1 border-0 bg-transparent py-1.5 pl-2 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
											/>
										</div>
									</div>
								</div>
							</div>
						</div>
						<div className="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8">
							<Link to="/location/create">
								<button type="button" className="text-sm font-semibold leading-6 text-gray-900">
									Cancel
								</button>
							</Link>
							<button
								type="submit"
								disabled={isPending}
								className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
							>
								{isPending ? "Adding..." : "Add"}
							</button>
						</div>
					</form>
				</div>
			</div>
		</>
	);
}
