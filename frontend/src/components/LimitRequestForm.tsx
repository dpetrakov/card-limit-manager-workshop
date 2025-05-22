import React, { useState } from "react";
import {
  LimitRequestCreate,
  LimitRequestView,
  ApiError,
} from "../types/limitRequest";
import { createLimitRequest } from "../services/limitRequestApi";
import styles from "./LimitRequestForm.module.css";

interface LimitRequestFormProps {
  onSuccess: (request: LimitRequestView) => void;
  onError: (error: string) => void;
}

const LimitRequestForm: React.FC<LimitRequestFormProps> = ({
  onSuccess,
  onError,
}) => {
  const [amount, setAmount] = useState<number | "">("");
  const [currency, setCurrency] = useState<string>("USD");
  const [justification, setJustification] = useState<string>("");
  const [desiredDate, setDesiredDate] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [formError, setFormError] = useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setFormError(null);
    onError(""); // Clear previous errors

    if (amount === "" || amount <= 0) {
      setFormError("Amount must be a positive number.");
      return;
    }
    if (!currency) {
      setFormError("Currency is required.");
      return;
    }
    if (!justification) {
      setFormError("Justification is required.");
      return;
    }
    if (!desiredDate) {
      setFormError("Desired date is required.");
      return;
    }

    const requestData: LimitRequestCreate = {
      amount: Number(amount),
      currency,
      justification,
      desired_date: desiredDate,
    };

    setIsLoading(true);
    try {
      const result = await createLimitRequest(requestData);
      onSuccess(result);
    } catch (error: any) {
      if (error instanceof Error) {
        onError(error.message);
      } else {
        onError("An unexpected error occurred.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      {formError && <div className={styles.formError}>{formError}</div>}

      <div className={styles.formGroup}>
        <label htmlFor="amount">Amount:</label>
        <input
          id="amount"
          type="number"
          value={amount}
          onChange={(e) =>
            setAmount(e.target.value === "" ? "" : parseFloat(e.target.value))
          }
          min="0.01"
          step="0.01"
          required
          className={styles.input}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="currency">Currency:</label>
        <input
          id="currency"
          type="text"
          value={currency}
          onChange={(e) => setCurrency(e.target.value.toUpperCase())}
          required
          maxLength={3}
          minLength={3}
          className={styles.input}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="justification">Justification:</label>
        <textarea
          id="justification"
          value={justification}
          onChange={(e) => setJustification(e.target.value)}
          required
          className={styles.textarea}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="desiredDate">Desired Date:</label>
        <input
          id="desiredDate"
          type="date"
          value={desiredDate}
          onChange={(e) => setDesiredDate(e.target.value)}
          required
          className={styles.input}
        />
      </div>

      <button type="submit" disabled={isLoading} className={styles.button}>
        {isLoading ? "Submitting..." : "Submit Request"}
      </button>
    </form>
  );
};

export default LimitRequestForm;
