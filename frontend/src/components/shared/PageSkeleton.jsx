import { Skeleton } from "../ui/skeleton";

export default function PageSkeleton({ cards = 4, rows = 6 }) {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: cards }).map((_, index) => (
          <Skeleton key={index} className="h-32 rounded-3xl" />
        ))}
      </div>
      <Skeleton className="h-96 rounded-3xl" />
      {Array.from({ length: rows }).map((_, index) => (
        <Skeleton key={index} className="h-14 rounded-2xl" />
      ))}
    </div>
  );
}
