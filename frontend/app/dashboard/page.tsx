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
  FileIcon,
  FolderIcon,
  DownloadIcon,
  EyeIcon,
  PlusIcon,
  ArrowRightLeftIcon,
  UploadIcon,
} from "@/components/icons";
import { MainLayout } from "@/components/layouts/main-layout";
import { FaDropbox, FaGoogleDrive } from "react-icons/fa6";
import { IconType } from "react-icons/lib";

export default function DashboardPage() {
  return (
    <MainLayout>
      <div className="flex flex-col gap-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight mb-2">
            Welcome back, Ashpak Veetar
          </h1>
          <p className="text-muted-foreground">
            {"Here's an overview of your cloud storage"}
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">
            Storage Overview
          </h2>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <StorageCard
              provider="Google Drive"
              Icon={FaGoogleDrive}
              used={15.4}
              total={30}
              color="#4285F4"
            />
            <StorageCard
              provider="Dropbox"
              Icon={FaDropbox}
              used={8.2}
              total={20}
              color="#0061FF"
            />
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg">Connect New Provider</CardTitle>
                <CardDescription>Add more cloud storage</CardDescription>
              </CardHeader>
              <CardContent className="flex justify-center items-center py-8">
                <Button variant="outline" className="h-12 px-6">
                  <PlusIcon className="mr-2 h-4 w-4" />
                  Connect Provider
                </Button>
              </CardContent>
            </Card>
          </div>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">Recent Files</h2>
          <Card>
            <CardContent className="p-0">
              <div className="rounded-md border">
                <div className="relative w-full overflow-auto">
                  <table className="w-full caption-bottom text-sm">
                    <thead>
                      <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Name
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Size
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Modified
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Provider
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Actions
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <RecentFileRow
                        name="Project Proposal.docx"
                        Icon={FileIcon}
                        size="2.4 MB"
                        modified="Today, 2:30 PM"
                        Provider={FaGoogleDrive}
                      />
                      <RecentFileRow
                        name="Financial Report.xlsx"
                        Icon={FileIcon}
                        size="4.8 MB"
                        modified="Yesterday"
                        Provider={FaDropbox}
                      />
                      <RecentFileRow
                        name="Marketing Assets"
                        Icon={FolderIcon}
                        size="128 MB"
                        modified="Aug 12, 2023"
                        Provider={FaGoogleDrive}
                      />
                      <RecentFileRow
                        name="Presentation.pptx"
                        Icon={FileIcon}
                        size="8.2 MB"
                        modified="Aug 10, 2023"
                        Provider={FaDropbox}
                      />
                      <RecentFileRow
                        name="Product Photos"
                        Icon={FolderIcon}
                        size="1.2 GB"
                        modified="Aug 5, 2023"
                        Provider={FaGoogleDrive}
                      />
                    </tbody>
                  </table>
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex justify-center p-4">
              <Button variant="outline">View All Files</Button>
            </CardFooter>
          </Card>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">
            Quick Actions
          </h2>
          <div className="grid gap-4 md:grid-cols-3">
            <Button className="h-auto py-6 flex flex-col gap-2">
              <PlusIcon className="h-6 w-6" />
              <span>Connect New Provider</span>
            </Button>
            <Button
              className="h-auto py-6 flex flex-col gap-2"
              variant="outline"
            >
              <ArrowRightLeftIcon className="h-6 w-6" />
              <span>Start New Transfer</span>
            </Button>
            <Button
              className="h-auto py-6 flex flex-col gap-2"
              variant="outline"
            >
              <UploadIcon className="h-6 w-6" />
              <span>Upload File</span>
            </Button>
          </div>
        </section>
      </div>
    </MainLayout>
  );
}

function StorageCard({
  provider,
  Icon,
  used,
  total,
  color,
}: {
  provider: string;
  Icon: IconType;
  used: number;
  total: number;
  color: string;
}) {
  const percentage = (used / total) * 100;

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{provider}</CardTitle>
          <Icon className="h-6 w-6" />
        </div>
        <CardDescription>Storage usage</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex justify-center">
          <div className="relative h-24 w-24">
            <svg className="h-full w-full" viewBox="0 0 100 100">
              <circle
                className="stroke-muted"
                cx="50"
                cy="50"
                r="40"
                fill="none"
                strokeWidth="10"
              />
              <circle
                className="transition-all duration-300 ease-in-out"
                cx="50"
                cy="50"
                r="40"
                fill="none"
                strokeWidth="10"
                stroke={color}
                strokeDasharray={`${percentage * 2.51} 251`}
                strokeDashoffset="0"
                strokeLinecap="round"
                transform="rotate(-90 50 50)"
              />
            </svg>
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-lg font-medium">
                {percentage.toFixed(0)}%
              </span>
            </div>
          </div>
        </div>
        <div className="text-center">
          <p className="text-sm text-muted-foreground">
            {used} GB of {total} GB used
          </p>
        </div>
      </CardContent>
      <CardFooter>
        <Button variant="outline" className="w-full">
          Browse Files
        </Button>
      </CardFooter>
    </Card>
  );
}

function RecentFileRow({
  name,
  Icon,
  size,
  modified,
  Provider,
}: {
  name: string;
  Icon: IconType;
  size: string;
  modified: string;
  Provider: IconType;
}) {
  return (
    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
      <td className="p-4 align-middle">
        <div className="flex items-center gap-2">
          <Icon className="h-6 w-6" />
          <span>{name}</span>
        </div>
      </td>
      <td className="p-4 align-middle">{size}</td>
      <td className="p-4 align-middle">{modified}</td>
      <td className="p-4 align-middle">
        <Provider className="h-6 w-6" />
      </td>
      <td className="p-4 align-middle">
        <div className="flex gap-2">
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <EyeIcon className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <DownloadIcon className="h-4 w-4" />
          </Button>
        </div>
      </td>
    </tr>
  );
}
