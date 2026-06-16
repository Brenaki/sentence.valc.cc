export type Quote = {
  id: number;
  quote: string;
  author: string;
  work?: string;
  categories: string[];
  like_quantity: number;
  deslike_quantity: number;
  created_at: string;
  updated_at: string;
};

// O backend retorna um único objeto em /quote-of-the-day.
export type QuoteOfTheDayResponse = Quote;
