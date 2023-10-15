import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import type { LocationReadResponse } from "../types";

export default function LocationCreateInMemory() {
	const navigate = useNavigate();

	const [todo, setTodo] = useState("");

	const submit = async () => {
		const payload = {
			kind: "memory",
		};
		const response = await fetch("/api/v1/locations", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(payload),
		});
		const data: LocationReadResponse = await response.json();
		navigate(`/locations/${data.location.id}`);
	};

	return (
		<div className="space-y-10 divide-y divide-gray-900/10">
			<div className="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
				<div className="px-4 sm:px-0">
					<h2 className="text-base font-semibold leading-7 text-gray-900">In-Memory Location</h2>
					<p className="mt-1 text-sm leading-6 text-gray-600">An in-memory location for testing Dripfile.</p>
				</div>

				<div className="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2">
					<div className="px-4 py-6 sm:p-8">
						<div className="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
							<div className="sm:col-span-4">
								<label htmlFor="todo" className="block text-sm font-medium leading-6 text-gray-900">
									TODO
								</label>
								<div className="mt-2">
									<div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
										<input
											type="text"
											name="todo"
											id="todo"
											value={todo}
											onChange={(e) => setTodo(e.target.value)}
											className="block flex-1 border-0 bg-transparent py-1.5 pl-2 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
										/>
									</div>
								</div>
							</div>
						</div>
					</div>
					<div className="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8">
						<Link to="/locations/create" type="button" className="text-sm font-semibold leading-6 text-gray-900">
							Cancel
						</Link>
						<button
							onClick={async (e) => {
								await submit();
							}}
							className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						>
							Add
						</button>
					</div>
				</div>
			</div>
		</div>
	);
}
