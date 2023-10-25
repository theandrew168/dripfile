export type ErrorResponse = {
	error: string | Record<string, string>;
};

export function isErrorResponse(response: any): response is ErrorResponse {
	return "error" in response;
}
