"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  AlertTriangleIcon,
  HomeIcon,
  ChevronDownIcon,
  XCircleIcon,
} from "lucide-react";

export default function ErrorPage() {
  const router = useRouter();
  const [showDetails, setShowDetails] = useState(false);

  // In a real app, this would be a real error code and message
  const errorCode = "ERR_SYSTEM_500";
  const errorDetails = "An unexpected error occurred in the application.";

  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh] text-center px-4">
      <div className="mb-6">
        <AlertTriangleIcon className="h-20 w-20 text-amber-500" />
      </div>

      <h1 className="text-3xl font-bold tracking-tight mb-2">
        Something went wrong
      </h1>
      <p className="text-muted-foreground max-w-md mb-4">
        {"We've encountered an unexpected error while processing your request."}
      </p>

      <div className="text-sm text-muted-foreground mb-8">
        Reference:{" "}
        <code className="bg-muted px-1 py-0.5 rounded">{errorCode}</code>
      </div>

      <Alert variant="destructive" className="max-w-md mb-8">
        <XCircleIcon className="h-4 w-4" />
        <AlertTitle>System Error</AlertTitle>
        <AlertDescription>
          {"The operation couldn't be completed. Please try again."}
        </AlertDescription>
      </Alert>

      <div className="w-full max-w-md space-y-6">
        <div className="flex flex-col sm:flex-row justify-center gap-4">
          <Button
            variant="outline"
            onClick={() => router.push("/dashboard")}
            className="flex items-center gap-2"
          >
            <HomeIcon className="h-4 w-4" />
            Go Home
          </Button>
        </div>

        <Collapsible
          open={showDetails}
          onOpenChange={setShowDetails}
          className="w-full"
        >
          <CollapsibleTrigger asChild>
            <Button variant="ghost" className="flex items-center gap-2 mx-auto">
              Technical Details
              <ChevronDownIcon
                className={`h-4 w-4 transition-transform ${showDetails ? "rotate-180" : ""}`}
              />
            </Button>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <div className="p-4 bg-muted rounded-md text-left mt-2">
              <p className="text-sm font-mono break-all">{errorDetails}</p>
            </div>
          </CollapsibleContent>
        </Collapsible>
      </div>
    </div>
  );
}
