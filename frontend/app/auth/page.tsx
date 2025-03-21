"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { useState } from "react";
import Link from "next/link";
import { CloudIcon } from "lucide-react";
import { FaDropbox } from "react-icons/fa6";
import { FcGoogle } from "react-icons/fc";
import { AppConfig } from "@/lib/config";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

const dropboxCodeSchema = z.object({
  code: z
    .string()
    .length(43, { message: "Please provide a valid dropbox code." }),
});

type DropboxCodeSchema = z.infer<typeof dropboxCodeSchema>;

export default function AuthPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <div className="flex justify-center mb-2">
            <CloudIcon className="h-12 w-12 text-primary" />
          </div>
          <CardTitle className="text-2xl font-bold">
            Sign in to Cloudmesh
          </CardTitle>
          <CardDescription>
            Access all your cloud files in one place
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <OAuthProviders />
          <div className="text-center text-sm text-muted-foreground">
            We use OAuth for secure authentication. No password required.
          </div>
        </CardContent>
        <CardFooter className="flex justify-center gap-4 text-sm text-muted-foreground">
          <Link href="#" className="hover:underline">
            Terms of Service
          </Link>
          <Link href="#" className="hover:underline">
            Privacy Policy
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}

function OAuthProviders() {
  const [showDropboxToken, setShowDropboxToken] = useState(false);

  const onSelectProvider = (provider: "google" | "dropbox") => {
    const { BACKEND_API_URL, BACKEND_API_VERSION } = AppConfig;
    const baseUrl = `${BACKEND_API_URL}/${BACKEND_API_VERSION}`;
    switch (provider) {
      case "dropbox":
        window.open(
          `${baseUrl}/api/auth/${provider}`,
          "_blank",
          "width=800,height=600",
        );
        break;
      default:
        window.location.href = `${baseUrl}/api/auth/${provider}`;
    }
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<DropboxCodeSchema>({ resolver: zodResolver(dropboxCodeSchema) });

  const onSubmitCode = ({ code }: DropboxCodeSchema) => {
    const { BACKEND_API_URL, BACKEND_API_VERSION } = AppConfig;
    const baseUrl = `${BACKEND_API_URL}/${BACKEND_API_VERSION}`;
    window.location.href = `${baseUrl}/api/auth/dropbox/callback?code=${code}`;
  };

  return (
    <div className="space-y-3">
      <Button
        variant="outline"
        className="w-full flex items-center gap-2 h-12"
        onClick={() => {
          onSelectProvider("google");
        }}
      >
        <FcGoogle className="h-6 w-6" />
        Continue with Google
      </Button>

      <Collapsible open={showDropboxToken}>
        <CollapsibleTrigger asChild>
          <Button
            variant="outline"
            className="w-full flex items-center gap-2 h-12"
            onClick={() => {
              if (!showDropboxToken) {
                onSelectProvider("dropbox");
                setShowDropboxToken(true);
              }
            }}
          >
            <FaDropbox color="blue" className="h-6 w-6" />
            Continue with Dropbox
          </Button>
        </CollapsibleTrigger>
        <CollapsibleContent className="mt-3 space-y-3">
          <div className="text-sm">
            Dropbox requires an additional step. Paste the token you received.
          </div>
          <form onSubmit={handleSubmit(onSubmitCode)}>
            <Input placeholder="Enter Dropbox token" {...register("code")} />
            <p className="text-red-500 text-sm">{errors.code?.message}</p>{" "}
            <Button type="submit" disabled={!isValid} className="w-full mt-2">
              Complete Connection
            </Button>
          </form>
        </CollapsibleContent>
      </Collapsible>
    </div>
  );
}
