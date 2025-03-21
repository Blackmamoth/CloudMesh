import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { FaDropbox, FaGoogleDrive } from "react-icons/fa6";
import { MainLayout } from "@/components/layouts/main-layout";
import { IconType } from "react-icons/lib";

export default function ConnectPage() {
  return (
    <MainLayout>
      <div className="flex flex-col gap-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight mb-2">
            Connect Cloud Providers
          </h1>
          <p className="text-muted-foreground">
            Connect your cloud storage accounts to access all your files in one
            place.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">
            Add New Provider
          </h2>
          <div className="grid gap-4 md:grid-cols-2">
            <ProviderCard
              name="Google Drive"
              Icon={FaGoogleDrive}
              description="Connect your Google Drive to access all your Google documents, spreadsheets, and other files."
            />

            <DropboxProviderCard />
          </div>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">
            Connected Providers
          </h2>
          <div className="space-y-4">
            <ConnectedProviderCard
              name="Google Drive"
              Icon={FaGoogleDrive}
              email="ashpakv88@gmail.com"
              connectedSince="March 8, 2025"
              used={7.3}
              total={15}
            />
            <ConnectedProviderCard
              name="Dropbox"
              Icon={FaDropbox}
              email="ashpakv88@gmail.com"
              connectedSince="March 10, 2025"
              used={0.7}
              total={2}
            />
          </div>
        </section>
      </div>
    </MainLayout>
  );
}

function ProviderCard({
  name,
  Icon,
  description,
}: {
  name: string;
  Icon: IconType;
  description: string;
}) {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-3">
          <Icon className="h-6 w-6" />
          <CardTitle>{name}</CardTitle>
        </div>
      </CardHeader>
      <CardContent>
        <CardDescription className="text-sm">{description}</CardDescription>
      </CardContent>
      <CardFooter>
        <Button className="w-full">Connected</Button>
      </CardFooter>
    </Card>
  );
}

function DropboxProviderCard() {
  return (
    <Dialog>
      <Card>
        <CardHeader>
          <div className="flex items-center gap-3">
            <FaDropbox className="h-8 w-8" />
            <CardTitle>Dropbox</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <CardDescription className="text-sm">
            Connect your Dropbox account to access all your Dropbox files and
            folders. Requires a two-step authentication process.
          </CardDescription>
        </CardContent>
        <CardFooter>
          <DialogTrigger asChild>
            <Button className="w-full">Connected</Button>
          </DialogTrigger>
        </CardFooter>
      </Card>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Complete Dropbox Connection</DialogTitle>
          <DialogDescription>
            Paste the token you received from Dropbox to complete the connection
            process.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <Input placeholder="Enter Dropbox token" />
          <p className="text-sm text-muted-foreground">
            This token is provided by Dropbox after the initial authorization
            step and is required to finalize the connection.
          </p>
        </div>
        <DialogFooter>
          <Button variant="outline">Cancel</Button>
          <Button>Complete Connection</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function ConnectedProviderCard({
  name,
  Icon,
  email,
  connectedSince,
  used,
  total,
}: {
  name: string;
  Icon: IconType;
  email: string;
  connectedSince: string;
  used: number;
  total: number;
}) {
  const percentage = (used / total) * 100;

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-3">
            <Icon className="h-6 w-6" />
            <div>
              <h3 className="font-medium">{name}</h3>
              <p className="text-sm text-muted-foreground">{email}</p>
            </div>
          </div>

          <div className="text-sm text-muted-foreground">
            Connected since {connectedSince}
          </div>

          <div className="w-full max-w-[200px]">
            <div className="flex justify-between text-sm mb-1">
              <span>{used} GB used</span>
              <span>{total} GB</span>
            </div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
              <div
                className="h-full bg-primary"
                style={{ width: `${percentage}%` }}
              />
            </div>
          </div>

          <div className="flex gap-2">
            <Button variant="outline" size="sm">
              Refresh
            </Button>
            <Button variant="destructive" size="sm">
              Disconnect
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
