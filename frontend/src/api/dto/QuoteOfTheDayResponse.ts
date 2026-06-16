export type Quote = {
  quote: string;
  author: string;
  work?: string;
  categories: string[];
};

export type QuoteOfTheDayResponse = Quote[];