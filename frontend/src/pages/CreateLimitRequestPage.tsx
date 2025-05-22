import React, { useState } from "react";
import LimitRequestForm from "../components/LimitRequestForm";
import { LimitRequestView } from "../types/limitRequest";
import styles from "../components/LimitRequestForm.module.css"; // Re-using styles for messages

const CreateLimitRequestPage: React.FC = () => {
  const [submissionResult, setSubmissionResult] =
    useState<LimitRequestView | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>("");

  const handleSuccess = (request: LimitRequestView) => {
    setSubmissionResult(request);
    setErrorMessage(""); // Clear any previous errors
  };

  const handleError = (error: string) => {
    setErrorMessage(error);
    setSubmissionResult(null); // Clear any previous success
  };

  return (
    <div>
      <h1>Create New Limit Request</h1>
      <LimitRequestForm onSuccess={handleSuccess} onError={handleError} />

      {errorMessage && (
        <div
          className={styles.errorMessage}
          style={{ maxWidth: "500px", margin: "1rem auto" }}
        >
          Error: {errorMessage}
        </div>
      )}

      {submissionResult && (
        <div
          className={styles.successMessage}
          style={{ maxWidth: "500px", margin: "1rem auto" }}
        >
          Request Submitted Successfully!
          <div className={styles.requestDetails}>
            <h3>Request Details:</h3>
            <p>
              <strong>ID:</strong> <code>{submissionResult.id}</code>
            </p>
            <p>
              <strong>Status:</strong> {submissionResult.status}
            </p>
            <p>
              <strong>Amount:</strong> {submissionResult.amount}{" "}
              {submissionResult.currency}
            </p>
            <p>
              <strong>Justification:</strong> {submissionResult.justification}
            </p>
            <p>
              <strong>Desired Date:</strong> {submissionResult.desired_date}
            </p>
            <p>
              <strong>User ID:</strong> <code>{submissionResult.user_id}</code>
            </p>
            {submissionResult.current_approver_id && (
              <p>
                <strong>Current Approver ID:</strong>{" "}
                <code>{submissionResult.current_approver_id}</code>
              </p>
            )}
            <p>
              <strong>Created At:</strong>{" "}
              {new Date(submissionResult.created_at).toLocaleString()}
            </p>
            <p>
              <strong>Updated At:</strong>{" "}
              {new Date(submissionResult.updated_at).toLocaleString()}
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default CreateLimitRequestPage;
