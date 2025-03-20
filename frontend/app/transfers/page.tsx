// eslint no-use-before-define: 0
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
// import { Progress } from "@/components/ui/progress";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  // ArrowRightIcon,
  /*  FileIcon, */ FolderIcon,
} from "@/components/icons";
import { MainLayout } from "@/components/layouts/main-layout";
import { FaGoogleDrive, FaDropbox } from "react-icons/fa6";
import { IconType } from "react-icons/lib";

export default function TransfersPage() {
  return (
    <MainLayout>
      <div className="flex flex-col gap-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight mb-2">
            Transfer Files
          </h1>
          <p className="text-muted-foreground">
            Move or copy files between your cloud storage providers
          </p>
        </section>

        <section>
          <Card>
            <CardHeader>
              <CardTitle>New Transfer</CardTitle>
              <CardDescription>
                Select source and destination for your file transfer
              </CardDescription>
            </CardHeader>
            <CardContent>
              <TransferForm />
            </CardContent>
          </Card>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold tracking-tight">
            Recent Transfers
          </h2>
          <div className="space-y-4">
            {/* <TransferHistoryItem
              status="completed"
              source="Google Drive"
              sourceIcon={<FaGoogleDrive className="h-4 w-4" />}
              destination="Dropbox"
              destinationIcon={<FaDropbox className="h-4 w-4" />}
              files={12}
              size="256 MB"
              startTime="Today, 2:30 PM"
              duration="5 minutes"
            />
            <TransferHistoryItem
              status="in-progress"
              source="Dropbox"
              sourceIcon={<FaDropbox className="h-4 w-4" />}
              destination="Google Drive"
              destinationIcon={<FaGoogleDrive className="h-4 w-4" />}
              files={128}
              size="1.2 GB"
              startTime="Today, 3:15 PM"
              progress={65}
            />
            <TransferHistoryItem
              status="failed"
              source="Google Drive"
              sourceIcon={<FaGoogleDrive className="h-4 w-4" />}
              destination="Dropbox"
              destinationIcon={<FaDropbox className="h-4 w-4" />}
              files={3}
              size="48 MB"
              startTime="Yesterday, 4:20 PM"
              error="Destination quota exceeded"
            /> */}
          </div>
        </section>
      </div>
    </MainLayout>
  );
}

function TransferForm() {
  const [transferType, setTransferType] = useState("copy");
  const [maintainStructure] = useState(true);

  return (
    <div className="space-y-8">
      <div className="grid gap-6 md:grid-cols-2">
        <div className="space-y-4">
          <div className="font-medium">Source</div>
          <Select defaultValue="google-drive">
            <SelectTrigger>
              <SelectValue placeholder="Select provider" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="google-drive">
                <div className="flex items-center gap-2">
                  <FaGoogleDrive className="h-4 w-4" />
                  <span>Google Drive</span>
                </div>
              </SelectItem>
              <SelectItem value="dropbox">
                <div className="flex items-center gap-2">
                  <FaDropbox className="h-4 w-4" />
                  <span>Dropbox</span>
                </div>
              </SelectItem>
            </SelectContent>
          </Select>

          <Card>
            <CardContent className="p-4 h-[300px] overflow-auto">
              <div className="space-y-2">
                <div className="flex items-center gap-2 p-2 rounded hover:bg-muted cursor-pointer">
                  <Checkbox id="source-select-all" />
                  <Label htmlFor="source-select-all">Select All</Label>
                </div>
                {/* <SourceFileItem
                  name="Documents"
                  icon={<FolderIcon className="h-4 w-4" />}
                  isFolder={true}
                />
                <SourceFileItem
                  name="Photos"
                  icon={<FolderIcon className="h-4 w-4" />}
                  isFolder={true}
                />
                <SourceFileItem
                  name="Project Proposal.docx"
                  icon={<FileIcon className="h-4 w-4" />}
                  size="2.4 MB"
                />
                <SourceFileItem
                  name="Financial Report.xlsx"
                  icon={<FileIcon className="h-4 w-4" />}
                  size="4.8 MB"
                />
                <SourceFileItem
                  name="Presentation.pptx"
                  icon={<FileIcon className="h-4 w-4" />}
                  size="8.2 MB"
                /> */}
              </div>
            </CardContent>
            <CardFooter className="border-t px-4 py-2">
              <div className="text-sm text-muted-foreground">
                3 items selected (15.4 MB)
              </div>
            </CardFooter>
          </Card>
        </div>

        <div className="space-y-4">
          <div className="font-medium">Destination</div>
          <Select defaultValue="dropbox">
            <SelectTrigger>
              <SelectValue placeholder="Select provider" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="google-drive">
                <div className="flex items-center gap-2">
                  <FaGoogleDrive className="h-4 w-4" />
                  <span>Google Drive</span>
                </div>
              </SelectItem>
              <SelectItem value="dropbox">
                <div className="flex items-center gap-2">
                  <FaDropbox className="h-4 w-4" />
                  <span>Dropbox</span>
                </div>
              </SelectItem>
            </SelectContent>
          </Select>

          <Card>
            <CardContent className="p-4 h-[300px] overflow-auto">
              <div className="space-y-2">
                <DestinationFolderItem
                  name="Root"
                  Icon={FolderIcon}
                  isSelected={true}
                />
                <DestinationFolderItem
                  name="Documents"
                  Icon={FolderIcon}
                  isSelected={false}
                />
                <DestinationFolderItem
                  name="Photos"
                  Icon={FolderIcon}
                  isSelected={false}
                />
                <DestinationFolderItem
                  name="Projects"
                  Icon={FolderIcon}
                  isSelected={false}
                />
              </div>
            </CardContent>
            <CardFooter className="border-t px-4 py-2 flex justify-between">
              <div className="text-sm text-muted-foreground">
                Selected: Root
              </div>
              <Button size="sm" variant="outline">
                New Folder
              </Button>
            </CardFooter>
          </Card>
        </div>
      </div>

      <div className="space-y-4">
        <div className="font-medium">Transfer Options</div>

        <div className="space-y-4">
          <div>
            <div className="mb-2">Operation</div>
            <RadioGroup
              defaultValue={transferType}
              onValueChange={setTransferType}
              className="flex flex-col space-y-1"
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="copy" id="copy" />
                <Label htmlFor="copy">Copy files (keep original files)</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="move" id="move" />
                <Label htmlFor="move">
                  Move files (delete original files after transfer)
                </Label>
              </div>
            </RadioGroup>
          </div>

          <div className="flex flex-col space-y-1">
            <div className="flex items-center space-x-2">
              <Checkbox
                id="maintain-structure"
                checked={maintainStructure}
                // onCheckedChange={setMaintainStructure}
              />
              <Label htmlFor="maintain-structure">
                Maintain folder structure
              </Label>
            </div>
          </div>

          <div>
            <div className="mb-2">Duplicate handling</div>
            <Select defaultValue="rename">
              <SelectTrigger className="w-full md:w-[250px]">
                <SelectValue placeholder="Select option" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="rename">Rename new files</SelectItem>
                <SelectItem value="replace">Replace existing files</SelectItem>
                <SelectItem value="skip">Skip duplicate files</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      <div className="flex justify-end">
        <Button size="lg">Start Transfer</Button>
      </div>
    </div>
  );
}

// function SourceFileItem({ name, icon, size, isFolder = false }) {
//   return (
//     <div className="flex items-center gap-2 p-2 rounded hover:bg-muted cursor-pointer">
//       <Checkbox id={`source-${name}`} />
//       <div className="flex items-center gap-2 flex-1">
//         {icon}
//         <Label htmlFor={`source-${name}`}>{name}</Label>
//       </div>
//       {!isFolder && <div className="text-xs text-muted-foreground">{size}</div>}
//     </div>
//   );
// }

function DestinationFolderItem({
  name,
  Icon,
  isSelected = false,
}: {
  name: string;
  Icon: IconType;
  isSelected: boolean;
}) {
  return (
    <div
      className={`flex items-center gap-2 p-2 rounded hover:bg-muted cursor-pointer ${isSelected ? "bg-muted" : ""}`}
    >
      <div className="flex items-center gap-2 flex-1">
        <Icon className="h-4 w-4" />
        <span>{name}</span>
      </div>
      {isSelected && (
        <div className="text-xs font-medium text-primary">Selected</div>
      )}
    </div>
  );
}

// function TransferHistoryItem({
//   status,
//   source,
//   sourceIcon,
//   destination,
//   destinationIcon,
//   files,
//   size,
//   startTime,
//   duration,
//   progress,
//   error,
// }) {
//   return (
//     <Card>
//       <CardContent className="p-4">
//         <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
//           <div className="flex items-center gap-2">
//             {status === "completed" && (
//               <div className="h-2 w-2 rounded-full bg-green-500"></div>
//             )}
//             {status === "in-progress" && (
//               <div className="h-2 w-2 rounded-full bg-blue-500"></div>
//             )}
//             {status === "failed" && (
//               <div className="h-2 w-2 rounded-full bg-red-500"></div>
//             )}
//             <div className="font-medium">
//               {status === "completed" && "Completed"}
//               {status === "in-progress" && "In Progress"}
//               {status === "failed" && "Failed"}
//             </div>
//           </div>

//           <div className="flex items-center gap-2">
//             <div className="flex items-center gap-1">
//               {sourceIcon}
//               <span>{source}</span>
//             </div>
//             <ArrowRightIcon className="h-4 w-4 text-muted-foreground" />
//             <div className="flex items-center gap-1">
//               {destinationIcon}
//               <span>{destination}</span>
//             </div>
//           </div>

//           <div className="text-sm">
//             {files} files ({size})
//           </div>

//           <div className="text-sm text-muted-foreground">
//             Started: {startTime}
//             {duration && <span> • Duration: {duration}</span>}
//           </div>

//           <div className="flex gap-2">
//             {status === "failed" && (
//               <Button variant="outline" size="sm">
//                 Retry
//               </Button>
//             )}
//             <Button variant="ghost" size="sm">
//               Details
//             </Button>
//           </div>
//         </div>

//         {status === "in-progress" && (
//           <div className="mt-4">
//             <div className="flex justify-between text-sm mb-1">
//               <span>{progress}%</span>
//               <span>Transferring...</span>
//             </div>
//             <Progress value={progress} />
//           </div>
//         )}

//         {status === "failed" && error && (
//           <div className="mt-4 text-sm text-red-500">Error: {error}</div>
//         )}
//       </CardContent>
//     </Card>
//   );
// }
