import { Card, CardContent } from "./ui/card";
import { Skeleton } from "./ui/skeleton";

export default function LoadingState({ rows = 4 }) {
  return (
    <Card>
      <CardContent className="space-y-4 p-6">
        <Skeleton className="h-6 w-40" />
        {Array.from({ length: rows }).map((_, index) => (
          <Skeleton key={index} className="h-16 w-full" />
        ))}
      </CardContent>
    </Card>
  );
}
