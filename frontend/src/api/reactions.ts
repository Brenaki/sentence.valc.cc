const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

// Valores aceitos pelo backend em POST /quotes/{id}/reactions.
export const Reaction = {
  Dislike: 0,
  Like: 1,
} as const;

export type Reaction = (typeof Reaction)[keyof typeof Reaction];

export async function reactToQuote(
  quoteId: number,
  reaction: Reaction,
): Promise<void> {
  const response = await fetch(`${API_URL}/quotes/${quoteId}/reactions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ reaction }),
  });

  if (response.status === 429) {
    throw new Error("Calma! Você está reagindo rápido demais. Tente de novo em instantes.");
  }

  if (!response.ok) {
    throw new Error(`Erro ao reagir (${response.status}): ${response.statusText}`);
  }
}
