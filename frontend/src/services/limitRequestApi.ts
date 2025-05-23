// Example for frontend/src/services/limitRequestApi.ts
import type {
  LimitRequestCreate,
  LimitRequestView,
  ApiError,
} from "../types/limitRequest";

const API_BASE_URL = "/api/v1"; // Proxied by Nginx

export async function createLimitRequest(
  requestData: LimitRequestCreate,
): Promise<LimitRequestView> {
  const response = await fetch(`${API_BASE_URL}/requests`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // Add Authorization header if/when auth is implemented
    },
    body: JSON.stringify(requestData),
  });

  if (!response.ok) {
    // Attempt to parse error response, but be prepared for non-JSON responses or network issues
    try {
      const errorData: ApiError = await response.json();
      throw new Error(
        errorData.message ||
          `API Error: ${response.status} ${response.statusText}`,
      );
    } catch (e) {
      // If response.json() fails or errorData.message is not present
      throw new Error(
        `API Error: ${response.status} ${response.statusText}. Failed to parse error response.`,
      );
    }
  }
  return response.json() as Promise<LimitRequestView>;
}
