"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { CloudIcon, HomeIcon, ArrowLeftIcon } from "lucide-react";

export default function NotFoundPage() {
  const router = useRouter();

  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh] text-center px-4">
      <div className="relative mb-8">
        <div className="text-[150px] font-bold text-primary/10 leading-none select-none">
          404
        </div>
        <div className="absolute inset-0 flex items-center justify-center"></div>
      </div>

      <h1 className="text-3xl font-bold tracking-tight mb-2">
        {"We've lost this cloud"}
      </h1>
      <p className="text-muted-foreground max-w-md mb-8">
        {
          "The page you're looking for doesn't exist or has been moved to a\
        different location."
        }
      </p>

      <div className="w-full max-w-md space-y-6">
        <div className="flex flex-col sm:flex-row justify-center gap-4">
          <Button
            variant="default"
            onClick={() => router.push("/dashboard")}
            className="flex items-center gap-2"
          >
            <HomeIcon className="h-4 w-4" />
            Go Home
          </Button>
          <Button
            variant="outline"
            onClick={() => router.back()}
            className="flex items-center gap-2"
          >
            <ArrowLeftIcon className="h-4 w-4" />
            Go Back
          </Button>
        </div>
      </div>

      <div className="mt-12 relative">
        <div className="absolute -top-6 left-1/2 transform -translate-x-1/2">
          <CloudIcon className="h-8 w-8 text-muted-foreground/30" />
        </div>
        <div className="absolute -top-10 left-1/4">
          <CloudIcon className="h-6 w-6 text-muted-foreground/20" />
        </div>
        <div className="absolute -top-8 right-1/4">
          <CloudIcon className="h-5 w-5 text-muted-foreground/20" />
        </div>
        <p className="text-sm text-muted-foreground pt-4 border-t max-w-md">
          If you believe this is an error, please contact our support team.
        </p>
      </div>
    </div>
  );
}
