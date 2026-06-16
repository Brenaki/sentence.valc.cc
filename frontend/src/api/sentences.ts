import type { Quote } from "./dto/QuoteOfTheDayResponse";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function fetchQuoteOfTheDay(): Promise<Quote> {
  const response = await fetch(`${API_URL}/quote-of-the-day`);

  if (!response.ok) {
    throw new Error(`Erro na API (${response.status}): ${response.statusText}`);
  }

  return (await response.json()) as Quote;
}
