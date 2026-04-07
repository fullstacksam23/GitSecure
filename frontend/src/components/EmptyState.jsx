import { SearchX } from "lucide-react";
import { Card, CardContent } from "./ui/card";

export default function EmptyState({ title, description, icon: Icon = SearchX }) {
  return (
    <Card className="border-dashed">
      <CardContent className="flex flex-col items-center justify-center px-8 py-14 text-center">
        <div className="mb-4 rounded-2xl bg-muted p-4">
          <Icon className="h-6 w-6 text-muted-foreground" />
        </div>
        <h3 className="font-display text-lg font-semibold text-slate-950">{title}</h3>
        <p className="mt-2 max-w-md text-sm leading-6 text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  );
}
