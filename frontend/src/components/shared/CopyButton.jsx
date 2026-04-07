import { Copy, CopyCheck } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "../ui/button";

export default function CopyButton({ value, label = "Copied" }) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    if (!value) return;
    await navigator.clipboard.writeText(value);
    setCopied(true);
    toast.success(label);
    window.setTimeout(() => setCopied(false), 1200);
  }

  return (
    <Button variant="ghost" size="icon" onClick={handleCopy} aria-label="Copy">
      {copied ? <CopyCheck className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
    </Button>
  );
}
