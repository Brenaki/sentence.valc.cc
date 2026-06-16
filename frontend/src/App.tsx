import "./App.css";
import { useEffect, useState } from "react";
import { Badge } from "./components/ui/badge";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card";
import { Skeleton } from "./components/ui/skeleton";
import { Spinner } from "./components/ui/spinner";
import { todayDate } from "./lib/functions/date";
import { fetchQuoteOfTheDay } from "./api/sentences";
import type { Quote } from "./api/dto/QuoteOfTheDayResponse";

function App() {
  const date = todayDate();
  const formattedDate = `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()}`;

  const [quote, setQuote] = useState<Quote | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;

    fetchQuoteOfTheDay()
      .then((data) => {
        if (active) setQuote(data);
      })
      .catch((err: Error) => {
        if (active) setError(err.message);
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, []);

  const searchUrl = quote?.work
    ? `https://www.google.com/search?q=${encodeURIComponent(quote.work)}`
    : null;

  return (
    <main className="flex items-center justify-center min-h-screen p-4 sm:p-6">
      <Card className="w-full max-w-sm sm:max-w-md lg:max-w-lg xl:max-w-xl animate-in fade-in zoom-in-95 slide-in-from-bottom-4 duration-700">
        <CardHeader className="flex justify-center">
          <CardTitle className="text-base sm:text-lg md:text-xl text-muted-foreground text-center">
            Sentence of Today - {formattedDate}
          </CardTitle>
        </CardHeader>

        <CardContent className="text-stone-800">
          {loading && (
            <div className="border-2 border-gray-800 rounded-xl p-4 sm:p-6 animate-in fade-in duration-500">
              <div className="flex items-center gap-2 text-stone-600 mb-4">
                <Spinner className="size-5" />
                <span className="text-sm sm:text-base font-medium">
                  Carregando frase...
                </span>
              </div>
              <div className="space-y-3">
                <Skeleton className="h-6 w-full rounded" />
                <Skeleton className="h-6 w-5/6 rounded" />
                <Skeleton className="h-4 w-1/3 rounded" />
              </div>
            </div>
          )}

          {error && !loading && (
            <div className="border-2 border-red-700 rounded-xl p-6 text-red-700 text-base animate-in fade-in duration-500">
              <p>Erro ao carregar frase: {error}</p>
            </div>
          )}

          {quote && !loading && (
            <div className="text-xl sm:text-2xl font-bold text-justify border-2 border-gray-800 rounded-xl p-4 sm:p-6 animate-in fade-in duration-500">
              <p>
                "{quote.quote}"{" "}
                <span className="text-base sm:text-lg text-stone-600">
                  - {quote.author}
                </span>
              </p>
            </div>
          )}
        </CardContent>

        {loading && (
          <CardFooter className="flex flex-col gap-6 sm:gap-8 justify-center animate-in fade-in duration-500">
            <div className="flex flex-wrap gap-2 sm:gap-3 justify-center">
              {Array.from({ length: 4 }).map((_, index) => (
                <Skeleton key={index} className="h-6 w-20 rounded-full" />
              ))}
            </div>
            <Skeleton className="h-4 w-2/3 rounded" />
          </CardFooter>
        )}

        {quote && !loading && (
          <CardFooter className="flex flex-col gap-6 sm:gap-8 justify-center animate-in fade-in slide-in-from-bottom-2 duration-700">
            <div className="flex flex-wrap gap-2 sm:gap-3 justify-center">
              {quote.categories.map((category, index) => (
                <Badge
                  key={index}
                  className="rounded-full transition-transform hover:scale-105"
                >
                  {category}
                </Badge>
              ))}
            </div>

            {quote.work && searchUrl && (
              <div className="flex text-center">
                <p className="text-sm sm:text-base">
                  Read in the original work:{" "}
                  <a
                    href={searchUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="underline text-blue-800 transition-colors hover:text-blue-600"
                  >
                    {quote.work}
                  </a>
                </p>
              </div>
            )}
          </CardFooter>
        )}
      </Card>
    </main>
  );
}

export default App;
