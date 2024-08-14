"use client";

import React, { useState } from "react";

export default function Page({ params }: { params: { id: string } }) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    setError(null);

    try {
      /*
      const response = await fetch("/api/submit", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Failed to submit the data. Please try again.");
      }

      const data = await response.json();
       */

      console.log(formData.get("media"));
    } catch (error) {
      // Capture the error message to display to the user
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <form action={onSubmit}>
        <input type="text" name="media" />
        <button type="submit" disabled={isLoading}>
          {isLoading ? "Loading..." : "送信"}
        </button>
      </form>
    </div>
  );
}
