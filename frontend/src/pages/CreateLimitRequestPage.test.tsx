import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom";
import CreateLimitRequestPage from "./CreateLimitRequestPage";
import { LimitRequestView } from "../types/limitRequest";

// Mock the LimitRequestForm component
jest.mock("../components/LimitRequestForm", () => {
  // Mock implementation of LimitRequestForm
  // It needs to call onSuccess or onError passed as props to simulate form behavior
  return jest.fn(({ onSuccess, onError }) => (
    <div data-testid="mock-limit-request-form">
      <button
        onClick={() =>
          onSuccess({
            id: "test-id-123",
            amount: 1500,
            currency: "EUR",
            justification: "Mocked success justification",
            desired_date: "2025-01-15",
            status: "APPROVED",
            user_id: "user-test-789",
            current_approver_id: null,
            created_at: new Date("2024-07-28T10:00:00Z").toISOString(),
            updated_at: new Date("2024-07-28T11:00:00Z").toISOString(),
          } as LimitRequestView)
        }
      >
        Simulate Success
      </button>
      <button onClick={() => onError("Mocked API error message")}>
        Simulate Error
      </button>
    </div>
  ));
});

describe("CreateLimitRequestPage", () => {
  test("renders the LimitRequestForm component", () => {
    render(<CreateLimitRequestPage />);
    expect(screen.getByTestId("mock-limit-request-form")).toBeInTheDocument();
    expect(screen.getByText("Create New Limit Request")).toBeInTheDocument();
  });

  test("displays success message and request details when form submission is successful", async () => {
    render(<CreateLimitRequestPage />);

    // Simulate the mocked form calling onSuccess
    fireEvent.click(screen.getByText("Simulate Success"));

    await waitFor(() => {
      expect(
        screen.getByText(/Request Submitted Successfully!/i),
      ).toBeInTheDocument();
    });

    // Check for displayed details (adjust based on what CreateLimitRequestPage displays)
    expect(screen.getByText(/ID:/i)).toHaveTextContent("ID: test-id-123");
    expect(screen.getByText(/Status:/i)).toHaveTextContent("Status: APPROVED");
    expect(screen.getByText(/Amount:/i)).toHaveTextContent("Amount: 1500 EUR");
    expect(screen.getByText(/Justification:/i)).toHaveTextContent(
      "Justification: Mocked success justification",
    );
    expect(screen.getByText(/Desired Date:/i)).toHaveTextContent(
      "Desired Date: 2025-01-15",
    );
    expect(screen.getByText(/User ID:/i)).toHaveTextContent(
      "User ID: user-test-789",
    );
    // current_approver_id is null in mock, so it shouldn't be rendered by the page logic
    expect(screen.queryByText(/Current Approver ID:/i)).not.toBeInTheDocument();
    expect(screen.getByText(/Created At:/i)).toBeInTheDocument(); // Check for presence, exact date match is tricky
    expect(screen.getByText(/Updated At:/i)).toBeInTheDocument(); // Check for presence
  });

  test("displays error message when form submission results in an error", async () => {
    render(<CreateLimitRequestPage />);

    // Simulate the mocked form calling onError
    fireEvent.click(screen.getByText("Simulate Error"));

    await waitFor(() => {
      expect(
        screen.getByText(/Error: Mocked API error message/i),
      ).toBeInTheDocument();
    });
    expect(
      screen.queryByText(/Request Submitted Successfully!/i),
    ).not.toBeInTheDocument();
  });
});
