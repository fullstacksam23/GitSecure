import { useMemo, useState } from "react";
import { LoaderCircle, Search, Sparkles } from "lucide-react";
import { Button } from "./ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Input } from "./ui/input";

function parseRepoInput(value) {
  const trimmed = value.trim().replace(/^https?:\/\/github\.com\//, "").replace(/\/$/, "");
  const [owner = "", repo = ""] = trimmed.split("/");
  console.log(owner, repo)
  return { owner, repo };
}

export default function ScanForm({ onSubmit, loading }) {
  const [repoInput, setRepoInput] = useState("");

  const parsed = useMemo(() => parseRepoInput(repoInput), [repoInput]);
  const isValid = parsed.owner && parsed.repo;

  function handleSubmit(event) {
    event.preventDefault();
    if (!isValid || loading) return;
    onSubmit(parsed, repoInput);
  }

  return (
    <Card className="overflow-hidden">
      <CardHeader className="bg-gradient-to-r from-sky-100 via-white to-amber-50">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-white shadow-sm">
          <Sparkles className="h-5 w-5 text-sky-700" />
        </div>
        <CardTitle>Start a repository scan</CardTitle>
        <CardDescription>
          Enter a GitHub repository in the form <span className="font-semibold">owner/repo</span>
          to kick off SBOM extraction and vulnerability analysis.
        </CardDescription>
      </CardHeader>
      <CardContent className="pt-6">
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-800" htmlFor="repo">
              GitHub repository
            </label>
            <div className="relative">
              <Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <Input
                id="repo"
                placeholder="e.g. vercel/next.js"
                value={repoInput}
                onChange={(event) => setRepoInput(event.target.value)}
                className="pl-11"
              />
            </div>
            <p className="text-sm text-muted-foreground">
              The current backend only supports public GitHub repositories.
            </p>
          </div>

          <Button className="w-full sm:w-auto" type="submit" disabled={!isValid || loading}>
            {loading ? <LoaderCircle className="h-4 w-4 animate-spin" /> : null}
            {loading ? "Starting scan..." : "Start scan"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
